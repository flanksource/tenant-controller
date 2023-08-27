package main

import (
	"fmt"
	"os"

	"github.com/flanksource/tenant-controller/cmd"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(commit) > 8 {
		version = fmt.Sprintf("%v, commit %v, built at %v", version, commit[0:8], date)
	}

	cmd.Root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version of tenant-controller",
		Args:  cobra.MinimumNArgs(0),
		Run: func(*cobra.Command, []string) {
			fmt.Println(version)
		},
	})

	if err := cmd.Root.Execute(); err != nil {
		os.Exit(1)
	}
}
