package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	cliName        = "cn"
	cliDescription = "Ceph Nano - One step S3 in container with Ceph."
)

var (
	// Version is the Ceph Nano version
	Version = "1.0.0"

	// WorkingDirectory is the working directory where objects can be put inside S3
	WorkingDirectory = "/usr/share/ceph-nano"

	// CephNanoUID is the uid of the S3 user
	CephNanoUID = "nano"

	// RgwPort is the rgw listenning port
	RgwPort = "8000"

	// ContainerName is name of the container
	ContainerName = "ceph-nano"

	rootCmd = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"cn"},
		//Long:
	}
)

// Main is the main function calling the whole program
func Main() {
	dockerExist()
	seLinux()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		CliStartNano(),
		CliStopNano(),
		CliStatusNano(),
		CliPurgeNano(),
		CliLogsNano(),
		CliRestartNano(),
		CliVersionNano(),
		cmdS3,
	)
}

func init() {
	cobra.EnableCommandSorting = false
}
