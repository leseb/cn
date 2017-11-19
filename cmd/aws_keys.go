package cmd

/*
import (
	"context"
	"fmt"

	"docker.io/go-docker/api/types"
	"github.com/docker/docker/client"
)
*/

// getAwsKey gets AWS keys from inside the container
func getAwsKey() {
	/*

		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}

		cmd := []string{"/bin/cat", "/nano_user_details"}

		options := types.ExecConfig{
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          cmd,
		}

		response, err := cli.ContainerExecCreate(ctx, ContainerName, options)
		if err != nil {
			panic(err)
		}

		// error: cannot use options (type types.ExecConfig) as type types.ExecStartCheck in argument to cli.ContainerExecAttach
		connection, err := cli.ContainerExecAttach(ctx, response.ID, options)
		if err != nil {
			panic(err)
		}

		defer connection.Close()
		output, err := connection.Reader.ReadString('\n')
		// time to parse json in go!
		fmt.Print(output)
	*/
}
