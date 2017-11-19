package cmd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// GrepForSuccess searchs for the word 'SUCCESS' inside the container logs
func GrepForSuccess() bool {
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

// CephNanoHealth loops on GrepForSuccess for 30 seconds, fails after.
func CephNanoHealth() {
	// remove that once the loop works
	GrepForSuccess()

	/*
		timeout := time.After(30 * time.Second)
		poll := 3 * time.Second
			for {
				Health := GrepForSuccess()
				select {
				case <-Health:
					fmt.Println("The end!")
					return
				case <-timeout:
					fmt.Println("There's no more time to this. Exiting!")
					return
				default:
					fmt.Println("still waiting")
				}
				time.Sleep(poll)
			}
	*/
}

// CephNanoS3Health loops for 30 seconds while testing Ceph RGW heatlh
func CephNanoS3Health() {
	ips, _ := getInterfaceIPv4s()
	/*
		Taking the first IP is probably not ideal
		IMHO, using the interface with most of the traffic is better
	*/
	var url string
	url = "http://" + ips[0].String() + ":8000"
	response, err := http.Get(url)
	if err != nil {
		panic(err)
		// decrement counter
	} else {
		defer response.Body.Close()
		_, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
	}
}
