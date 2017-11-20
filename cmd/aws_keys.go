package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// getAwsKey gets AWS keys from inside the container
func getAwsKey() (string, string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	//cmd := []string{"/bin/cat", "/nano_user_details"}
	cmd := []string{"cat", "/nano_user_details"}

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
	//output, err := ioutil.ReadAll(connection.Reader)
	output, err := ioutil.ReadFile("/tmp/lol")
	if err != nil {
		panic(err)
	}

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
	if err != nil {
		panic(err)
	}

	var CephNanoAccessKey string
	CephNanoAccessKey = parsedMap.Keys[0].AccessKey
	var CephNanoSecretKey string
	CephNanoSecretKey = parsedMap.Keys[0].SecretKey
	return CephNanoAccessKey, CephNanoSecretKey

	//return "a", "b"
}
