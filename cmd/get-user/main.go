package main

import (
	"github.com/KnoblauchPilze/go-game/pkg/dtos"
	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get-user",
	Short: "Get details about an existing user",
	Args:  cobra.RangeArgs(0, 1),
	Run:   getUserCmdBody,
}

func main() {
	logger.Configure(logger.Configuration{
		Service: "get-user",
		Level:   logrus.DebugLevel,
	})

	if err := getCmd.Execute(); err != nil {
		logger.Fatalf("get-user command failed (err: %v)", err)
	}
}

func getUserCmdBody(cmd *cobra.Command, args []string) {
	ud := dtos.UserDto{
		Mail:     "toto@some-mail.com",
		Name:     "toto",
		Password: "123456",
	}

	if len(args) > 0 {
		ud.Mail = args[0]
	}

	logger.Infof("getting details for %+v", ud)

	if err := doServerRequest(ud); err != nil {
		logger.Errorf("get operation failed (err: %v)", err)
	}
}

func doServerRequest(in dtos.UserDto) error {
	return errors.NewCode(errors.ErrNotImplemented)
}
