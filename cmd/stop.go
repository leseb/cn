package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

// stopNano stops the container
func stopNano(c *cli.Context) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	timeout := 5 * time.Second
	if status := containerStatus(false); !status {
		fmt.Println("ceph-nano is already stopped!")
		os.Exit(1)
	} else {
		fmt.Println("Stopping ceph-nano... ")
		if err := cli.ContainerStop(ctx, ContainerName, &timeout); err != nil {
			panic(err)
		}
	}
}
