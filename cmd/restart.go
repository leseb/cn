package cmd

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

func RestartNano(c *cli.Context) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	fmt.Println("Restarting ceph-nano \n")
	if err := cli.ContainerRestart(ctx, "ceph-nano", nil); err != nil {
		panic(err)
	}
}
