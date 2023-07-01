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

var deleteCmd = &cobra.Command{
	Use:   "delete-user",
	Short: "Delete an existing user",
	Args:  cobra.ExactArgs(1),
	Run:   deleteUserCmdBody,
}

const serverUrl = "http://localhost:3000"

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
	id := args[0]

	logger.Infof("deleting user %s", id)

	if err := doServerRequest(id); err != nil {
		logger.Errorf("delete operation failed (err: %v)", err)
	}
}

func doServerRequest(id string) error {
	url := fmt.Sprintf("%s/users/%s", serverUrl, id)

	rb := connection.NewHttpDeleteRequestBuilder()
	rb.SetUrl(url)

	req, err := rb.Build()
	if err != nil {
		return err
	}
	resp, err := req.Perform()
	if err != nil {
		return err
	}

	var out dtos.DeleteResponse
	if err = rest.GetBodyFromHttpResponseAs(resp, &out); err != nil {
		return err
	}

	logger.Infof("server response: %+v", out)

	return nil
}
