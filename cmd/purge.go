package cmd

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/urfave/cli"
)

func removeContainer() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	options := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         true,
	}
	// we don't necessarily want to catch errors here
	// it's not an issue if the container does not exist
	cli.ContainerRemove(ctx, ContainerName, options)
}

// purgeNano purges Ceph Nano.
func purgeNano(c *cli.Context) {
	removeContainer()
	removeCephNanoVolumes()
}
