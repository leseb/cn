package cmd

// CreateCephNanoVolumes creates docker volumes for the container
func CreateCephNanoVolumes() {
	/*
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}

		CephDockerVols := []string{"varlibceph", "etcceph"}

		for _, vol := range CephDockerVols {
			volume, err := cli.VolumeCreate(ctx, volumetypes.VolumesCreateBody{
				Name:   vol,
				Driver: "local",
			})
			if err != nil {
				panic(err)
			}
		}
	*/
}

// RemoveCephNanoVolumes removes docker volumes
func RemoveCephNanoVolumes() {
	/*
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}

		CephDockerVols := []string{"varlibceph", "etcceph"}

		for _, vol := range CephDockerVols {
			volume, err := cli.VolumeRemove(ctx, volumetypes.VolumesCreateBody{
				Name:   vol,
				Driver: "local",
			})
			if err != nil {
				panic(err)
			}
		}
	*/
}
