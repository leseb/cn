package cmd

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// CliStartNano is the Cobra CLI call
func CliStartNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts object storage server",
		Args:  cobra.NoArgs,
		Run:   startNano,
		Example: "cn start --work-dir /tmp \n" +
			"cn start",
	}
	cmd.Flags().StringVarP(&WorkingDirectory, "work-dir", "d", "/usr/share/ceph-nano", "Directory to work from")
	return cmd
}

// startNano starts Ceph Nano
func startNano(cmd *cobra.Command, args []string) {
	// Test for a leftover container
	// Usually happens when someone fails to run the container on an exposed directory
	// Typical error on Docker For Mac you will see:
	// panic: Error response from daemon: Mounts denied:
	// The path /usr/share/ceph-nano is not shared from OS X and is not known to Docker.
	// You can configure shared paths from Docker -> Preferences... -> File Sharing.
	if status := containerStatus(true, "created"); status {
		removeContainer(ContainerName)
	}

	if status := containerStatus(false, "running"); status {
		fmt.Println("ceph-nano is already running!")
	} else if status := containerStatus(true, "exited"); status {
		fmt.Println("Starting ceph-nano...")
		startContainer()
	} else {
		fmt.Println("Running ceph-nano...")
		runContainer()
	}
	echoInfo()
}

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

	// TODO --memory 512m

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
