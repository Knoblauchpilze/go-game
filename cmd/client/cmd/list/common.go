package list

import (
	"errors"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/auth"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var userArg string
var tokenArg string

func Commands() []*cobra.Command {
	return []*cobra.Command{usersCmd, userCmd}
}

func buildTokenFromFlags() (auth.Token, error) {
	var token auth.Token
	var err error

	token.User, err = uuid.Parse(userArg)
	if err != nil {
		return token, err
	}

	if len(tokenArg) == 0 {
		return token, errors.New("invalid user token provided")
	}

	token.Value = tokenArg
	token.Expiration = time.Now().Add(1 * time.Minute)

	return token, nil
}
