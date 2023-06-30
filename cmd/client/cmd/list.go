package cmd

import (
	"github.com/KnoblauchPilze/go-game/cmd/client/cmd/list"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List information in the server",
}

func init() {
	listCmds := list.Commands()
	listCmd.AddCommand(listCmds...)

	rootCmd.AddCommand(listCmd)
}
