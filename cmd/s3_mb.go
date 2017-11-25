package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CliS3CmdMb is the Cobra CLI call
func CliS3CmdMb() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mb [BUCKET]",
		Short: "Make bucket",
		Args:  cobra.ExactArgs(1),
		Run:   S3CmdMb,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdMb wraps s3cmd command in the container
func S3CmdMb(cmd *cobra.Command, args []string) {
	if status := containerStatus(false, "running"); !status {
		fmt.Println("ceph-nano does not exist yet!")
		os.Exit(1)
	}
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is not running!")
		os.Exit(1)
	}
	command := []string{"s3cmd", "mb", "s3://" + args[0]}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
