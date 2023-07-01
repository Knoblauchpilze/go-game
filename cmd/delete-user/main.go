package main

import (
	"github.com/KnoblauchPilze/go-game/pkg/dtos"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete-user",
	Short: "Delete an existing user",
	Args:  cobra.RangeArgs(0, 1),
	Run:   deleteUserCmdBody,
}

func main() {
	logger.Configure(logger.Configuration{
		Service: "delete-user",
		Level:   logrus.DebugLevel,
	})

	if err := deleteCmd.Execute(); err != nil {
		logger.Fatalf("delete-user command failed (err: %v)", err)
	}
}

func deleteUserCmdBody(cmd *cobra.Command, args []string) {
	ud := dtos.UserDto{
		Mail:     "toto@some-mail.com",
		Name:     "toto",
		Password: "123456",
	}

	if len(args) > 0 {
		ud.Mail = args[0]
	}

	logger.Infof("deleting data for %+v", ud)

	if err := doServerRequest(ud); err != nil {
		logger.Errorf("delete operation failed (err: %v)", err)
	}
}

func doServerRequest(in dtos.UserDto) error {
	return errors.NewCode(errors.ErrNotImplemented)
}
