package main

import (
	"bytes"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"log"
	"os"
	"os/exec"
	"strings"
	//	"time"
)

var Version = "1.0.0"
var WorkingDirectory = "/usr/share/ceph-nano"
var CephNanoUid = "nano"
var RgwPort = "8000"
var CephNanoAccessKey = "accesskey"
var CephNanoSecretKey = "secretkey"
var ContainerName = "ceph-nano"

// Help template for ceph-nano.
/*
var nanoHelpTemplate = `NAME:
  {{.Name}} - {{.Usage}}

USAGE:
  {{.HelpName}} {{if .VisibleFlags}}[FLAGS] {{end}}COMMAND{{if .VisibleFlags}}{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}

COMMANDS:
  {{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}
FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}

VERSION:
  ` + Version +
	`{{ "\n"}}`
*/

func DockerExist() {
	if _, err := os.Stat("/usr/local/bin/docker"); os.IsNotExist(err) {
		fmt.Print("Docker is not present on your system, installation instructions can be found at https://docs.docker.com/engine/installation/")
		os.Exit(1)
	}
}

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

func GetAwsKey() {
	/*
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}

		cmd := []string{"/bin/cat", "/nano_user_details"}

		options := types.ExecConfig{
			AttachStdout: true,
			AttachStderr: true,
			Cmd:          cmd,
		}

			response, err := cli.ContainerExecCreate(ctx, ContainerName, options)
			if err != nil {
				panic(err)
			}

					// error: cannot use options (type types.ExecConfig) as type types.ExecStartCheck in argument to cli.ContainerExecAttach
					connection, err := cli.ContainerExecAttach(ctx, response.ID, options)
					if err != nil {
						panic(err)
					}

				defer connection.Close()

				output, err := connection.Reader.ReadString('\n')

				fmt.Print(output)
	*/
}

func GrepForSuccess() bool {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	out, err := cli.ContainerLogs(ctx, ContainerName, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	newStr := buf.String()

	if strings.Contains(newStr, "SUCCESS") {
		return true
	}
	return false
}

func CephNanoHealth() {
	// remove that once the loop works
	GrepForSuccess()

	/*
		timeout := time.After(30 * time.Second)
		poll := 3 * time.Second
			for {
				Health := GrepForSuccess()
				select {
				case <-Health:
					fmt.Println("The end!")
					return
				case <-timeout:
					fmt.Println("There's no more time to this. Exiting!")
					return
				default:
					fmt.Println("still waiting")
				}
				time.Sleep(poll)
			}
	*/
}

func CephNanoS3Health() {
	// curl --fail --silent --output /dev/null http://"$IP":8000
}

func WaitForContainer() {
	CephNanoHealth()
	CephNanoS3Health()
}

func Selinux() {
	if _, err := os.Stat("/sbin/getenforce"); !os.IsNotExist(err) {
		out, err := exec.Command("getenforce").Output()
		if err != nil {
			log.Fatal(err)
		}
		if string(out) == "Enforcing" {
			exec.Command("chcon -Rt svirt_sandbox_file_t %s", WorkingDirectory)
		}
	}
}

func StartNano(c *cli.Context) {
	if _, err := os.Stat(WorkingDirectory); os.IsNotExist(err) {
		os.Mkdir(WorkingDirectory, 0755)
	}
	CreateCephNanoVolumes()

	fmt.Println("Starting ceph-nano...")
	if status, _ := ContainerStatus(); status {
		StartContainer()
	} else {
		RunContainer()
	}

	// wait for the container to be ready
	//WaitForContainer()
}

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

func EchoInfo() {
	WaitForContainer()
	ContainerStatus()
	InfoLine := "\n" +
		"Ceph status is: $(docker exec ceph-nano ceph health) \n" +
		"Ceph Rados Gateway address is: http://$IP:$RGW_CIVETWEB_PORT \n" +
		"Your working directory is: " +
		WorkingDirectory +
		"\n" +
		"S3 user is: nano \n" +
		"S3 access key is: " +
		CephNanoAccessKey +
		"\n" +
		"S3 secret key is: " +
		CephNanoSecretKey +
		"\n" +
		""
	fmt.Println(InfoLine)
	GetAwsKey()
}

func ContainerStatus() (bool, string) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/ceph-nano" && container.State == "running" {
				fmt.Print("ceph-nano exists and running")
				return true, container.State
			}
		}
	}
	return false, "stopped"
}

func StopNano(c *cli.Context) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	fmt.Println("Stopping ceph-nano \n")
	if err := cli.ContainerStop(ctx, ContainerName, nil); err != nil {
		panic(err)
	}
}

func RestartNano(c *cli.Context) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	fmt.Println("Restarting ceph-nano \n")
	if err := cli.ContainerRestart(ctx, "ceph-nano", nil); err != nil {
		panic(err)
	}
}

func StatusNano(c *cli.Context) {
	EchoInfo()
}

func PurgeNano(c *cli.Context) {
}

func LogsNano(c *cli.Context) {
}

func S3cmdWrapper() {
}

func main() {
	app := cli.NewApp()
	app.UsageText = "ceph-nano [FLAGS] COMMAND [arguments...]"
	app.Name = "ceph-nano"
	app.Author = "ceph.com"
	app.Usage = "One step S3 in container with Ceph!"
	//	app.CustomAppHelpTemplate = nanoHelpTemplate

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "work-dir, d",
			Value: "" + WorkingDirectory,
			Usage: "Only files within this `DIRECTORY` can be uploaded in Ceph Nano.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "Starts object storage server. Default working directory is " + WorkingDirectory,
			Action: StartNano,
		},
		{
			Name:   "stop",
			Usage:  "Stops object storage server.",
			Action: StopNano,
		},
		{
			Name:   "restart",
			Usage:  "Restarts object storage server.",
			Action: RestartNano,
		},
		{
			Name:   "status",
			Usage:  "Prints useful information about the object storage server.",
			Action: StatusNano,
		},
		{
			Name:   "purge",
			Usage:  "DANGEROUS, removes the object storage server and all its data",
			Action: PurgeNano,
		},
		{
			Name:   "logs",
			Usage:  "Displays container S3 logs (can be really verbose)",
			Action: LogsNano,
		},
		{
			Name:      "mb",
			Usage:     "Make bucket",
			ArgsUsage: "[BUCKET]",
			Category:  "INTERACT WITH S3",
		},
		{
			Name:      "rb",
			Usage:     "Remove bucket",
			ArgsUsage: "[BUCKET]",
			Category:  "INTERACT WITH S3",
		},
		{
			Name:      "ls",
			Usage:     "List objects or buckets",
			ArgsUsage: "[BUCKET]",
			Category:  "INTERACT WITH S3",
		},
		{
			Name:     "la",
			Usage:    "List all object in all buckets",
			Category: "INTERACT WITH S3",
		},
		{
			Name:      "put",
			Usage:     "Put file into bucket",
			ArgsUsage: "FILE [BUCKET]",
			Category:  "INTERACT WITH S3",
		},
		{
			Name:      "get",
			Usage:     "Get file from bucket",
			ArgsUsage: "BUCKET/OBJECT LOCAL_FILE",
			Category:  "INTERACT WITH S3",
		},
		{
			Name:      "del",
			Usage:     "Delete file from bucket",
			ArgsUsage: "[BUCKET]/OBJECT",
			Category:  "INTERACT WITH S3",
		},
		{
			Name:      "sync",
			Usage:     "Synchronize a directory tree to S3",
			ArgsUsage: "LOCAL_DIR BUCKET[/PREFIX]",
			Category:  "INTERACT WITH S3",
		},
		{
			Name:      "du",
			Usage:     "Disk usage by buckets",
			ArgsUsage: "[BUCKET[/PREFIX]]",
			Category:  "INTERACT WITH S3",
		},
		{
			Name:      "info",
			Usage:     "Get various information about Buckets or Files",
			ArgsUsage: "BUCKET[/OBJECT]",
			Category:  "INTERACT WITH S3",
		},
		{
			Name:      "cp",
			Usage:     "Copy object",
			ArgsUsage: "BUCKET1/OBJECT1 BUCKET2[/OBJECT2]",
			Category:  "INTERACT WITH S3",
		},
		{
			Name:      "mv",
			Usage:     "Move object",
			ArgsUsage: "BUCKET1/OBJECT1 BUCKET2[/OBJECT2]",
			Category:  "INTERACT WITH S3",
		},
	}

	DockerExist()
	app.Run(os.Args)
}
