package main

import (
	"github.com/KnoblauchPilze/go-game/cmd/client/cmd"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
)

// https://github.com/spf13/cobra/blob/main/user_guide.md
func main() {
	logger.Configure(logger.Configuration{
		Service: "client",
	})

	cmd.Execute()
}
