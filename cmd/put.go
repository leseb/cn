package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CliS3CmdPut is the Cobra CLI call
func CliS3CmdPut() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put FILE [BUCKET]",
		Short: "Put file into bucket",
		Args:  cobra.ExactArgs(2),
		Run:   S3CmdPut,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

// S3CmdPut wraps s3cmd command in the container
func S3CmdPut(cmd *cobra.Command, args []string) {
	if status := containerStatus(true, "exited"); status {
		fmt.Println("ceph-nano is not running!")
		os.Exit(1)
	}
	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		fmt.Printf("Input file: '%s' doesn't exit in the current directory. \n"+
			"Use the full path of the file or change directory to your working directory: %s \n", args[0], WorkingDirectory)
		os.Exit(1)
	}
	command := []string{"s3cmd", "put", args[0], "s3://" + args[1]}
	output := execContainer(ContainerName, command)
	fmt.Printf("%s", output)
}
