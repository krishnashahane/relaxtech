package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BuildVersion holds the current release tag, injected at build time.
var BuildVersion = "0.1.0-dev"

var buildVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current relaxtech version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(BuildVersion)
	},
}
