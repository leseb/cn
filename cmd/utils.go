package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// validateEnv verifies the ability to run the program
func validateEnv() {
	dockerAPIVersion()
	dockerExist()
	seLinux()
}

// dockerApiVersion checks docker's API Version
func dockerAPIVersion() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	sv := fmt.Sprint(cli.ServerVersion(ctx))
	if err != nil {
		panic(err)
	}

	if strings.Contains(sv, "is too new") {
		ss := strings.SplitAfter(sv, "Maximum supported API version is ")
		os.Setenv("DOCKER_API_VERSION", ss[1])
	} else if strings.Contains(sv, "client is newer than server") {
		ss := strings.SplitAfter(sv, "server API version: ")
		// trim last character since this 'ss[1]' is '1.24.'
		os.Setenv("DOCKER_API_VERSION", ss[1][:len(ss[1])-1])
	}
}

// dockerExist makes sure Docker is installed
func dockerExist() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	_, err = cli.Info(ctx)
	if err != nil {
		fmt.Println("Docker is not present on your system or not started.\n" +
			"Make sure it's started or follow installation instructions at https://docs.docker.com/engine/installation/")
		os.Exit(1)
	}
}

// seLinux checks if Selinux is installed and set to Enforcing,
// we relabel our WorkingDirectory to allow the container to access files in this directory
func seLinux() {
	if _, err := os.Stat("/sbin/getenforce"); !os.IsNotExist(err) {
		out, err := exec.Command("getenforce").Output()
		if err != nil {
			panic(err)
		}
		if string(out) == "Enforcing" {
			if _, err := os.Stat(WorkingDirectory); os.IsNotExist(err) {
				os.Mkdir(WorkingDirectory, 0755)
			}
			exec.Command("sudo chcon -Rt svirt_sandbox_file_t %s", WorkingDirectory)
		}
	}
}

// byLastOctetValue implements sort.Interface used in sorting a list
// of ip address by their last octet value.
type byLastOctetValue []net.IP

func (n byLastOctetValue) Len() int      { return len(n) }
func (n byLastOctetValue) Swap(i, j int) { n[i], n[j] = n[j], n[i] }
func (n byLastOctetValue) Less(i, j int) bool {
	return []byte(n[i].To4())[3] < []byte(n[j].To4())[3]
}

// getInterfaceIPv4s is synonymous to net.InterfaceAddrs()
// returns net.IP IPv4 only representation of the net.Addr.
// Additionally the returned list is sorted by their last
// octet value.
//
// [The logic to sort by last octet is implemented to
// prefer CIDRs with higher octects, this in-turn skips the
// localhost/loopback address to be not preferred as the
// first ip on the list. Subsequently this list helps us print
// a user friendly message with appropriate values].
func getInterfaceIPv4s() ([]net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("Unable to determine network interface address. %s", err)
	}
	// Go through each return network address and collate IPv4 addresses.
	var nips []net.IP
	for _, addr := range addrs {
		if addr.Network() == "ip+net" {
			var nip net.IP
			// Attempt to parse the addr through CIDR.
			nip, _, err = net.ParseCIDR(addr.String())
			if err != nil {
				return nil, fmt.Errorf("Unable to parse addrss %s, error %s", addr, err)
			}
			// Collect only IPv4 addrs.
			if nip.To4() != nil {
				nips = append(nips, nip)
			}
		}
	}
	// Sort the list of IPs by their last octet value.
	sort.Sort(sort.Reverse(byLastOctetValue(nips)))
	return nips, nil
}

// execContainer execs a given command inside the container
func execContainer(ContainerName string, cmd []string) []byte {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	optionsCreate := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	response, err := cli.ContainerExecCreate(ctx, ContainerName, optionsCreate)
	if err != nil {
		panic(err)
	}

	optionsAttach := types.ExecStartCheck{
		Detach: false,
		Tty:    false,
	}
	connection, err := cli.ContainerExecAttach(ctx, response.ID, optionsAttach)
	if err != nil {
		panic(err)
	}

	defer connection.Close()
	output, err := ioutil.ReadAll(connection.Reader)
	if err != nil {
		panic(err)
	}

	// Remove 8 first characters to get a readable content
	// Sometimes the command returns nothing, without the following if the program fails without
	// runtime error: slice bounds out of range
	if len(output) > 0 {
		return output[8:]
	}
	return nil
}

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

// echoInfo prints useful information about Ceph Nano
func echoInfo() {
	// Always wait the container to be ready
	cephNanoHealth()
	cephNanoS3Health()

	// Fetch Amazon Keys
	CephNanoAccessKey, CephNanoSecretKey := getAwsKey()

	// Get Ceph health
	cmd := []string{"ceph", "health"}
	c := execContainer(ContainerName, cmd)

	// Get IPs, later using the first IP of the list is not ideal
	// However, Docker binds RGW port on 0.0.0.0 so any address will work
	ips, _ := getInterfaceIPv4s()

	// Get the working directory
	dir := dockerInspect()

	InfoLine := "\n" +
		strings.TrimSpace(string(c)) +
		" is the Ceph status. \n" +
		"S3 object server address is: http://" + ips[0].String() + ":8000 \n" +
		"S3 user is: nano \n" +
		"S3 access key is: " +
		CephNanoAccessKey +
		"\n" +
		"S3 secret key is: " +
		CephNanoSecretKey +
		"\n" +
		"Your working directory is: " +
		dir +
		"\n" +
		""
	fmt.Println(InfoLine)
}

// getAwsKey gets AWS keys from inside the container
func getAwsKey() (string, string) {
	cmd := []string{"/bin/cat", "/nano_user_details"}

	output := execContainer(ContainerName, cmd)

	// declare structures for json
	type s3Details []struct {
		AccessKey string `json:"Access_key"`
		SecretKey string `json:"Secret_key"`
	}
	type jason struct {
		Keys s3Details
	}
	// assign variable to our json struct
	var parsedMap jason

	json.Unmarshal(output, &parsedMap)

	var CephNanoAccessKey string
	CephNanoAccessKey = parsedMap.Keys[0].AccessKey
	var CephNanoSecretKey string
	CephNanoSecretKey = parsedMap.Keys[0].SecretKey
	return CephNanoAccessKey, CephNanoSecretKey
}

// dockerInspect inspect the container Binds
func dockerInspect() string {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	inspect, err := cli.ContainerInspect(ctx, ContainerName)
	if err != nil {
		panic(err)
	}
	parts := strings.Split(inspect.HostConfig.Binds[0], ":")
	return parts[0]
}

// inspectImage inspect a given image
func inspectImage() map[string]string {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	i, _, err := cli.ImageInspectWithRaw(ctx, ImageName)
	if err != nil {
		panic(err)
	}
	return i.Config.Labels
}

// pullImage downloads the container image
func pullImage() bool {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	_, _, err = cli.ImageInspectWithRaw(ctx, ImageName)
	if err != nil {
		fmt.Print("The container image is not present, pulling it. \n" +
			"This operation can take a few minutes.")

		out, err := cli.ImagePull(ctx, ImageName, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}

		reader := bufio.NewReader(out)
		defer out.Close() // pullResp is io.ReadCloser
		var respo bytes.Buffer
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				// it could be EOF or read error
				break
			}
			respo.Write(line)
			respo.WriteByte('\n')
			fmt.Print(".")
		}
		fmt.Println("")
		return true
	}
	return false
}
