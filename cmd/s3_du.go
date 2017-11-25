package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CliS3CmdDu is the Cobra CLI call
func CliS3CmdDu() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "du [BUCKET[/PREFIX]]",
		Short: "Disk usage by buckets",
		Args:  cobra.ExactArgs(1),
		Run:   S3CmdDu,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdDu wraps s3cmd command in the container
func S3CmdDu(cmd *cobra.Command, args []string) {
	if status := containerStatus(false, "running"); !status {
		fmt.Println("ceph-nano does not exist yet!")
		os.Exit(1)
	}
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is not running!")
		os.Exit(1)
	}
	command := []string{"s3cmd", "du", "s3://" + args[0]}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
