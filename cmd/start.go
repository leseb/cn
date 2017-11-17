package cmd

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"os"
)

// This function creates a new container when nothing exists
func RunContainer() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	imageName := "ceph/daemon"

	_, err = cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        imageName,
		ExposedPorts: nat.PortSet{"8000": struct{}{}},
	}, &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{nat.Port("8000"): {{HostIP: "127.0.0.1", HostPort: "8000"}}},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println(resp.ID)

}

// This function starts a container that is stopped
func StartContainer() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, ContainerName, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func StartNano(c *cli.Context) {
	if _, err := os.Stat(WorkingDirectory); os.IsNotExist(err) {
		os.Mkdir(WorkingDirectory, 0755)
	}
	CreateCephNanoVolumes()

	fmt.Println("Starting ceph-nano...")
	StartContainer()
	if status := ContainerStatus(); status {
		StartContainer()
	} else {
		RunContainer()
	}
	// wait for the container to be ready
	//WaitForContainer()
}
