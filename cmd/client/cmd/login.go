package cmd

import (
	"github.com/KnoblauchPilze/go-game/cmd/client/session"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/KnoblauchPilze/go-game/pkg/types"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the server with the specified credentials",
	Args:  cobra.RangeArgs(0, 2),
	Run:   loginCmdBody,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func loginCmdBody(cmd *cobra.Command, args []string) {
	ud := types.UserData{
		Name:     "toto",
		Password: "123456",
	}

	if len(args) > 0 {
		ud.Name = args[0]
	}
	if len(args) > 1 {
		ud.Password = args[1]
	}

	logger.Infof("logging in for %+v", ud)

	sess := session.NewManager(defaultServerURL)
	if err := sess.Login(ud); err != nil {
		logger.Errorf("failed to log in (%v)", err)
		return
	}
}
