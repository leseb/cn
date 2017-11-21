package cmd

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

// restartNano restarts the container
func restartNano(c *cli.Context) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	fmt.Println("Restarting ceph-nano...")
	if err := cli.ContainerRestart(ctx, "ceph-nano", nil); err != nil {
		panic(err)
	}
	echoInfo()
}
