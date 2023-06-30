package cmd

import (
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "Client allows to interact with the server (duh)",
	Long:  "Makes submitting commands to the server easy and fast",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatalf("failed to execute root command (%v)", err)
	}
}
