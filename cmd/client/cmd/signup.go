package cmd

import (
	"github.com/KnoblauchPilze/go-game/cmd/client/session"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/KnoblauchPilze/go-game/pkg/types"
	"github.com/spf13/cobra"
)

var signUpCmd = &cobra.Command{
	Use:   "sign-up",
	Short: "Sign up to the server with the specified user",
	Args:  cobra.RangeArgs(0, 2),
	Run:   signUpCmdBody,
}

func init() {
	rootCmd.AddCommand(signUpCmd)
}

func signUpCmdBody(cmd *cobra.Command, args []string) {
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

	logger.Infof("signing up for %+v", ud)

	sess := session.NewManager(defaultServerURL)
	if err := sess.SignUp(ud); err != nil {
		logger.Errorf("failed to sign up (%v)", err)
		return
	}
}
