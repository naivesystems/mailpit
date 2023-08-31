package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/axllent/mailpit/config"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the current version information",
	Long:  `Display the current version information (if available).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("%s %s compiled with %s on %s/%s\n",
			os.Args[0], config.Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
