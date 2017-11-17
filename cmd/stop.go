package cmd

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

func StopNano(c *cli.Context) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	fmt.Println("Stopping ceph-nano... \n")
	if err := cli.ContainerStop(ctx, ContainerName, nil); err != nil {
		panic(err)
	}
}
