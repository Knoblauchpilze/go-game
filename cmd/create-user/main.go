package main

import (
	"fmt"

	"github.com/KnoblauchPilze/go-game/pkg/connection"
	"github.com/KnoblauchPilze/go-game/pkg/dtos"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/KnoblauchPilze/go-game/pkg/rest"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create a new user",
	Args:  cobra.RangeArgs(0, 3),
	Run:   createUserCmdBody,
}

const serverUrl = "http://localhost:3000"

func main() {
	logger.Configure(logger.Configuration{
		Service: "create-user",
		Level:   logrus.DebugLevel,
	})

	if err := createCmd.Execute(); err != nil {
		logger.Fatalf("create-user command failed (err: %v)", err)
	}
}

func createUserCmdBody(cmd *cobra.Command, args []string) {
	ud := dtos.UserDto{
		Mail:     "toto@some-mail.com",
		Name:     "toto",
		Password: "123456",
	}

	if len(args) > 0 {
		ud.Mail = args[0]
	}
	if len(args) > 1 {
		ud.Name = args[1]
	}
	if len(args) > 2 {
		ud.Password = args[2]
	}

	logger.Infof("creating new user %+v", ud)

	if err := doServerRequest(ud); err != nil {
		logger.Errorf("create operation failed (err: %v)", err)
	}
}

func doServerRequest(in dtos.UserDto) error {
	var out dtos.SignUpResponse

	url := fmt.Sprintf("%s/users", serverUrl)

	rb := connection.NewHttpPostRequestBuilder()
	rb.SetUrl(url)
	rb.SetBody("application/json", in)

	req, err := rb.Build()
	if err != nil {
		return err
	}
	resp, err := req.Perform()
	if err != nil {
		return err
	}

	err = rest.GetBodyFromHttpResponseAs(resp, &out)
	if err != nil {
		return err
	}

	logger.Infof("server response: %+v", out)

	return nil
}
