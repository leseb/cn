package cmd

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

// runContainer creates a new container when nothing exists
func runContainer() {
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

	//CephDockerVols := []string{"varlibceph", "etcceph"}

	//TODO: change 0.0.0.0?
	exposedPorts, portBindings, _ := nat.ParsePortSpecs([]string{":8000:8000"})
	envs := []string{
		"DEBUG=verbose",
		"CEPH_DEMO_UID=" + CephNanoUID,
		"NETWORK_AUTO_DETECT=4",
		"MON_IP=127.0.0.1",
		"RGW_CIVETWEB_PORT=" + RgwPort,
		"CEPH_DAEMON=demo"}

	config := &container.Config{
		Image:        imageName,
		Hostname:     ContainerName + "-faa32aebf00b",
		ExposedPorts: exposedPorts,
		Env:          envs,
		Volumes: map[string]struct{}{
			"/etc/ceph":     struct{}{},
			"/var/lib/ceph": struct{}{},
		},
	}
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Binds:        []string{WorkingDirectory + ":/tmp"},
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, ContainerName)
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

// startContainer starts a container that is stopped
func startContainer() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, ContainerName, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

// startNano starts Ceph Nano
func startNano(c *cli.Context) {
	if _, err := os.Stat(WorkingDirectory); os.IsNotExist(err) {
		os.Mkdir(WorkingDirectory, 0755)
	}

	// test for leftover container
	if status := containerStatus(true, "created"); status {
		removeContainer(ContainerName)
	}

	if status := containerStatus(false, "running"); status {
		fmt.Println("ceph-nano is already running!")
		echoInfo()
	} else if status := containerStatus(true, "exited"); status {
		fmt.Println("Starting ceph-nano...")
		startContainer()
		// wait for the container to be ready
		echoInfo()
	} else {
		fmt.Println("Running ceph-nano...")
		runContainer()
		// wait for the container to be ready
		echoInfo()
	}
}
