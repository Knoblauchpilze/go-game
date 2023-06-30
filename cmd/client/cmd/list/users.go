package list

import (
	"github.com/KnoblauchPilze/go-game/cmd/client/session"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "List all registered users",
	Run:   usersCmdBody,
}

func init() {
	usersCmd.Flags().StringVar(&userArg, "user", "", "the id of the user")
	usersCmd.Flags().StringVar(&tokenArg, "token", "", "the token of the user")
}

func usersCmdBody(cmd *cobra.Command, args []string) {
	token, err := buildTokenFromFlags()
	if err != nil {
		logger.Errorf("invalid parameters to list users (%v)", err)
		return
	}

	sess := session.NewManager(defaultServerURL)
	if err := sess.Authenticate(token); err != nil {
		logger.Fatalf("Failed to list users: %+v", err)
		return
	}

	data, err := sess.ListUsers()
	if err != nil {
		logger.Fatalf("Failed to list users: %+v", err)
		return
	}

	logger.Infof("Users: %+v", data)
}
