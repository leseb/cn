package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func cephHealth() string {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	//cmd := []string{"/bin/cat", "/nano_user_details"}
	cmd := []string{"ceph", "health"}

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
	return string(output)
}

// echoInfo prints useful information about Ceph Nano
func echoInfo() {
	// Always wait the container to be ready
	cephNanoHealth()
	cephNanoS3Health()

	// Fetch Amazon Keys
	CephNanoAccessKey, CephNanoSecretKey := getAwsKey()

	// Get Ceph health
	c := cephHealth()

	// Get IPs
	ips, _ := getInterfaceIPv4s()

	InfoLine := strings.TrimSpace(c) +
		" is the Ceph status. \n" +
		"Ceph Rados Gateway address is: http://" + ips[0].String() + ":8000 \n" +
		"Your working directory is: " +
		WorkingDirectory +
		"\n" +
		"S3 user is: nano \n" +
		"S3 access key is: " +
		CephNanoAccessKey +
		"\n" +
		"S3 secret key is: " +
		CephNanoSecretKey +
		"\n" +
		""
	fmt.Println(InfoLine)
}
