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
// the parameter corresponds to the type listOptions and its entry all
func containerStatus(allList bool, containerState string) bool {
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

	// run the loop on both indexes, it's fine they have the same length
	for _, container := range containers {
		for i := range container.Names {
			if container.Names[i] == "/ceph-nano" && container.State == containerState {
				return true
			}
		}
	}
	return false
}

// statusNano show Ceph Nano status
func statusNano(c *cli.Context) {
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is stopped!")
		os.Exit(1)
	}
	echoInfo()
}
