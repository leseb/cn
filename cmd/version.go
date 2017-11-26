package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CliVersionNano is the Cobra CLI call
func CliVersionNano() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Ceph Nano",
		Args:  cobra.NoArgs,
		Run:   versionNano,
	}
	return cmd
}

// versionNano print Ceph Nano version
func versionNano(cmd *cobra.Command, args []string) {
	fmt.Println("Ceph Nano version: " + Version)
	ii := inspectImage()
	fmt.Println("Ceph Nano container image version: " + ii["head"])

}
