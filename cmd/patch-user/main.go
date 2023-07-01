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

var patchCmd = &cobra.Command{
	Use:   "patch-user",
	Short: "Patch data of an existing user",
	Args:  cobra.RangeArgs(0, 3),
	Run:   patchUserCmdBody,
}

const serverUrl = "http://localhost:3000"

func main() {
	logger.Configure(logger.Configuration{
		Service: "patch-user",
		Level:   logrus.DebugLevel,
	})

	if err := patchCmd.Execute(); err != nil {
		logger.Fatalf("patch-user command failed (err: %v)", err)
	}
}

func patchUserCmdBody(cmd *cobra.Command, args []string) {
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

	logger.Infof("updating details for %+v", ud)

	if err := doServerRequest(ud); err != nil {
		logger.Errorf("update operation failed (err: %v)", err)
	}
}

func doServerRequest(in dtos.UserDto) error {
	url := fmt.Sprintf("%s/users/%s", serverUrl, in.Id)

	rb := connection.NewHttpPatchRequestBuilder()
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

	var out dtos.PatchResponse
	if err = rest.GetBodyFromHttpResponseAs(resp, &out); err != nil {
		return err
	}

	logger.Infof("server response: %+v", out)

	return nil
}
