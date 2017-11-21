package cmd

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/urfave/cli"
)

func removeContainer(name string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	options := types.ContainerRemoveOptions{
		RemoveLinks:   false,
		RemoveVolumes: true,
		Force:         true,
	}
	// we don't necessarily want to catch errors here
	// it's not an issue if the container does not exist
	cli.ContainerRemove(ctx, name, options)
}

// purgeNano purges Ceph Nano.
func purgeNano(c *cli.Context) {
	fmt.Println("Purging ceph-nano... ")
	removeContainer(ContainerName)
}
