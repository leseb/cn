package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// grepForSuccess searchs for the word 'SUCCESS' inside the container logs
func grepForSuccess() bool {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	out, err := cli.ContainerLogs(ctx, ContainerName, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	newStr := buf.String()

	if strings.Contains(newStr, "SUCCESS") {
		return true
	}
	return false
}

// cephNanoHealth loops on grepForSuccess for 30 seconds, fails after.
func cephNanoHealth() {
	// setting timeout values
	var timeout int
	timeout = 60
	var poll int
	poll = 0

	// setting docker connection
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// wait for 60sec to validate that the container started properly
	for poll < timeout {
		if health := grepForSuccess(); health {
			return
		}
		time.Sleep(time.Second * 1)
		poll++
	}

	// if we reach here, something is broken in the container
	fmt.Print("The container from Ceph Nano never reached a clean state. Show the container logs:")
	// ideally we would return the second value of GrepForSuccess when it's false
	// this would mean having 2 return values for GrepForSuccess
	out, err := cli.ContainerLogs(ctx, ContainerName, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	newStr := buf.String()
	fmt.Println(newStr)
	fmt.Println("Please open an issue at: https://github.com/ceph/ceph-container.")
	os.Exit(1)
}

// curlS3 queries S3 URL
func curlS3() bool {
	ips, _ := getInterfaceIPv4s()
	// Taking the first IP is probably not ideal
	// IMHO, using the interface with most of the traffic is better
	var url string
	url = "http://" + ips[0].String() + ":8000"

	response, err := http.Get(url)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	if _, err := ioutil.ReadAll(response.Body); err != nil {
		return false
	}
	return true
}

// CephNanoS3Health loops for 30 seconds while testing Ceph RGW heatlh
func cephNanoS3Health() {
	// setting timeout
	var timeout int
	timeout = 30
	var poll int
	poll = 0

	for poll < timeout {
		if s3Health := curlS3(); s3Health {
			return
		}
		time.Sleep(time.Second * 1)
		poll++
	}
	fmt.Println("S3 gateway is not responding. Showing S3 logs:")
	showS3Logs()
	fmt.Println("Please open an issue at: https://github.com/ceph/ceph-nano.")
	os.Exit(1)
}
