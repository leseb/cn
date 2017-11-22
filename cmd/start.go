package cmd

import (
	"fmt"
	"os"
	"strings"

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
		runContainer(cmd, args)
	}
	echoInfo()
}

// runContainer creates a new container when nothing exists
func runContainer(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	imageName := "ceph/daemon"
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()

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

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	/* The if removes the error:
	 panic: runtime error: invalid memory address or nil pointer dereference
	[signal SIGSEGV: segmentation violation code=0x1 addr=0x20 pc=0x137a2b4]

	goroutine 1 [running]:
	github.com/leseb/cn/cmd.runContainer(0xc420089200, 0xc4202d8140, 0x0, 0x2)
	        /Users/leseb/go/src/github.com/leseb/cn/cmd/start.go:102 +0x734
	github.com/leseb/cn/cmd.startNano(0xc420089200, 0xc4202d8140, 0x0, 0x2)
	        /Users/leseb/go/src/github.com/leseb/cn/cmd/start.go:49 +0x1c3
	github.com/spf13/cobra.(*Command).execute(0xc420089200, 0xc4202d8120, 0x2, 0x2, 0xc420089200, 0xc4202d8120)
	        /Users/leseb/go/src/github.com/spf13/cobra/command.go:704 +0x2c6
	github.com/spf13/cobra.(*Command).ExecuteC(0x16a0520, 0xc4202dc0f0, 0xc4202c86c0, 0xc4202c8900)
	        /Users/leseb/go/src/github.com/spf13/cobra/command.go:785 +0x30e
	github.com/spf13/cobra.(*Command).Execute(0x16a0520, 0x0, 0xc4202c9440)
	        /Users/leseb/go/src/github.com/spf13/cobra/command.go:738 +0x2b
	github.com/leseb/cn/cmd.Main()
	        /Users/leseb/go/src/github.com/leseb/cn/cmd/main.go:47 +0x3b
	main.main()
	        /Users/leseb/go/src/github.com/leseb/cn/main.go:30 +0x20
	*/
	if err != nil {
		if strings.Contains(err.Error(), "Mounts denied") {
			fmt.Println("ERROR: It looks like you need to use the --work-dir option.")
			cmd.Help()
			os.Exit(1)
		} else {
			panic(err)
		}
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
