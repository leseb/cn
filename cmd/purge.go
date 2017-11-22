package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// CliPurgeNano is the Cobra CLI call
func CliPurgeNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge",
		Short: "DANGEROUS! Purges object storage server",
		Args:  cobra.NoArgs,
		Run:   purgeNano,
	}
	return cmd
}

// purgeNano purges Ceph Nano.
func purgeNano(cmd *cobra.Command, args []string) {
	if status := containerStatus(false, "running"); !status {
		fmt.Println("ceph-nano does not exist yet!")
		os.Exit(1)
	}
	fmt.Println("Purging ceph-nano... ")
	removeContainer(ContainerName)
}

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
