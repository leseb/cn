package cmd

import (
	"os"

	"github.com/urfave/cli"
)

// NanoCli is the CLI function
func nanoCli() {
	app := cli.NewApp()
	app.UsageText = "ceph-nano [FLAGS] COMMAND [arguments...]"
	app.Name = "ceph-nano"
	app.Author = "ceph.com"
	app.Usage = "One step S3 in container with Ceph!"
	app.Version = Version

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
			Action: startNano,
		},
		{
			Name:   "stop",
			Usage:  "Stops object storage server.",
			Action: stopNano,
		},
		{
			Name:   "restart",
			Usage:  "Restarts object storage server.",
			Action: restartNano,
		},
		{
			Name:   "status",
			Usage:  "Prints useful information about the object storage server.",
			Action: statusNano,
		},
		{
			Name:   "purge",
			Usage:  "DANGEROUS, removes the object storage server and all its data",
			Action: purgeNano,
		},
		{
			Name:   "logs",
			Usage:  "Displays container S3 logs (can be really verbose)",
			Action: logsNano,
		},
		{
			Name:      "mb",
			Usage:     "Make bucket",
			ArgsUsage: "[BUCKET]",
			Category:  "INTERACT WITH S3",
			Action:    s3cmdWrapper,
		},
		{
			Name:      "rb",
			Usage:     "Remove bucket",
			ArgsUsage: "[BUCKET]",
			Category:  "INTERACT WITH S3",
			Action:    s3cmdWrapper,
		},
		{
			Name:      "ls",
			Usage:     "List objects or buckets",
			ArgsUsage: "[BUCKET]",
			Category:  "INTERACT WITH S3",
			Action:    s3cmdWrapper,
		},
		{
			Name:     "la",
			Usage:    "List all object in all buckets",
			Category: "INTERACT WITH S3",
			Action:   s3cmdWrapper,
		},
		{
			Name:      "put",
			Usage:     "Put file into bucket",
			ArgsUsage: "FILE [BUCKET]",
			Category:  "INTERACT WITH S3",
			Action:    s3cmdWrapper,
		},
		{
			Name:      "get",
			Usage:     "Get file from bucket",
			ArgsUsage: "BUCKET/OBJECT LOCAL_FILE",
			Category:  "INTERACT WITH S3",
			Action:    s3cmdWrapper,
		},
		{
			Name:      "del",
			Usage:     "Delete file from bucket",
			ArgsUsage: "[BUCKET]/OBJECT",
			Category:  "INTERACT WITH S3",
			Action:    s3cmdWrapper,
		},
		{
			Name:      "sync",
			Usage:     "Synchronize a directory tree to S3",
			ArgsUsage: "LOCAL_DIR BUCKET[/PREFIX]",
			Category:  "INTERACT WITH S3",
			Action:    s3cmdWrapper,
		},
		{
			Name:      "du",
			Usage:     "Disk usage by buckets",
			ArgsUsage: "[BUCKET[/PREFIX]]",
			Category:  "INTERACT WITH S3",
			Action:    s3cmdWrapper,
		},
		{
			Name:      "info",
			Usage:     "Get various information about Buckets or Files",
			ArgsUsage: "BUCKET[/OBJECT]",
			Category:  "INTERACT WITH S3",
			Action:    s3cmdWrapper,
		},
		{
			Name:      "cp",
			Usage:     "Copy object",
			ArgsUsage: "BUCKET1/OBJECT1 BUCKET2[/OBJECT2]",
			Category:  "INTERACT WITH S3",
			Action:    s3cmdWrapper,
		},
		{
			Name:      "mv",
			Usage:     "Move object",
			ArgsUsage: "BUCKET1/OBJECT1 BUCKET2[/OBJECT2]",
			Category:  "INTERACT WITH S3",
			Action:    s3cmdWrapper,
		},
	}
	app.Run(os.Args)
}
