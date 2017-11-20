package cmd

import (
	"context"

	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

// createCephNanoVolumes creates docker volumes for the container
func createCephNanoVolumes() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	CephDockerVols := []string{"varlibceph", "etcceph"}
	for _, vol := range CephDockerVols {
		options := volume.VolumesCreateBody{
			Name:   vol,
			Driver: "local",
		}
		_, err = cli.VolumeCreate(ctx, options)
		if err != nil {
			panic(err)
		}
	}
}

// removeCephNanoVolumes removes docker volumes
func removeCephNanoVolumes() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	CephDockerVols := []string{"varlibceph", "etcceph"}
	for _, vol := range CephDockerVols {
		err = cli.VolumeRemove(ctx, vol, true)
		if err != nil {
			panic(err)
		}
	}

}
