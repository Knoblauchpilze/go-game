package cmd

import (
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/spf13/cobra"
)

var versionMajor = 0
var versionMinor = 1
var versionPatch = 0

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Client",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Infof("Client v%d.%d.%d", versionMajor, versionMinor, versionPatch)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
