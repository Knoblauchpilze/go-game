package main

import (
	"github.com/KnoblauchPilze/go-game/pkg/dtos"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update-user",
	Short: "Update data of an existing user",
	Args:  cobra.RangeArgs(0, 2),
	Run:   updateUserCmdBody,
}

func main() {
	logger.Configure(logger.Configuration{
		Service: "update-user",
		Level:   logrus.DebugLevel,
	})

	if err := updateCmd.Execute(); err != nil {
		logger.Fatalf("update-user command failed (err: %v)", err)
	}
}

func updateUserCmdBody(cmd *cobra.Command, args []string) {
	ud := dtos.UserDto{
		Mail:     "toto@some-mail.com",
		Name:     "toto",
		Password: "123456",
	}

	if len(args) > 0 {
		ud.Name = args[0]
	}
	if len(args) > 1 {
		ud.Password = args[1]
	}

	logger.Infof("updating details for %+v", ud)

	if err := doServerRequest(ud); err != nil {
		logger.Errorf("update operation failed (err: %v)", err)
	}
}

func doServerRequest(in dtos.UserDto) error {
	return errors.NewCode(errors.ErrNotImplemented)
}
