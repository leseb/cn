package cmd

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

// containerStatus checks container status
func containerStatus() bool {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/ceph-nano" && container.State == "running" {
				return true
			}
		}
	}
	fmt.Println("ceph-nano is stopped!")
	os.Exit(1)
	return false
}

// statusNano show Ceph Nano status
func statusNano(c *cli.Context) {
	containerStatus()
	echoInfo()
}
