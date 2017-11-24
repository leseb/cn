package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CliS3CmdLa is the Cobra CLI call
func CliS3CmdLa() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "la",
		Short: "List all object in all buckets",
		Args:  cobra.NoArgs,
		Run:   S3CmdLa,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdLa wraps s3cmd command in the container
func S3CmdLa(cmd *cobra.Command, args []string) {
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is not running!")
		os.Exit(1)
	}
	command := []string{"s3cmd", "la"}
	output := execContainer(ContainerName, command)
	if len(output) == 1 {
		command := []string{"s3cmd", "ls"}
		output := execContainer(ContainerName, command)
		fmt.Printf("%s", output)
	} else {
		fmt.Printf("%s", output)
	}
}
