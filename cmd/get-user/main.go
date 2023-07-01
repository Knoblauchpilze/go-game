package main

import (
	"fmt"
	"net/http"

	"github.com/KnoblauchPilze/go-game/pkg/connection"
	"github.com/KnoblauchPilze/go-game/pkg/dtos"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/KnoblauchPilze/go-game/pkg/rest"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get-user",
	Short: "Get details about an existing user",
	Args:  cobra.RangeArgs(0, 1),
	Run:   getUserCmdBody,
}

const serverUrl = "http://localhost:3000"

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
	var err error

	if len(args) == 0 {
		logger.Infof("get all users")
		err = getAllUsers()
	} else {
		logger.Infof("get details for user %s", args[0])
		err = getUser(args[0])
	}

	if err != nil {
		logger.Errorf("get operation failed (err: %v)", err)
	}
}

func getAllUsers() error {
	url := fmt.Sprintf("%s/users", serverUrl)

	resp, err := doServerRequest(url)
	if err != nil {
		return err
	}

	var out dtos.GetAllResponse
	if err = rest.GetBodyFromHttpResponseAs(resp, &out); err != nil {
		return err
	}

	logger.Infof("server response: %+v", out)

	return nil
}

func getUser(id string) error {
	url := fmt.Sprintf("%s/users/%s", serverUrl, id)

	resp, err := doServerRequest(url)
	if err != nil {
		return err
	}

	var out dtos.GetResponse
	if err = rest.GetBodyFromHttpResponseAs(resp, &out); err != nil {
		return err
	}

	logger.Infof("server response: %+v", out)

	return nil
}

func doServerRequest(url string) (*http.Response, error) {
	rb := connection.NewHttpGetRequestBuilder()
	rb.SetUrl(url)

	req, err := rb.Build()
	if err != nil {
		return nil, err
	}
	resp, err := req.Perform()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
