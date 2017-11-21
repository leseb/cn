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
func containerStatus(allList bool) bool {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	listOptions := types.ContainerListOptions{
		All:   allList,
		Quiet: true,
	}
	containers, err := cli.ContainerList(context.Background(), listOptions)
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/ceph-nano" {
				return true
			}
		}
	}
	return false
}

// statusNano show Ceph Nano status
func statusNano(c *cli.Context) {
	if status := containerStatus(false); !status {
		fmt.Println("ceph-nano is stopped!")
		os.Exit(1)
	}
	echoInfo()
}
