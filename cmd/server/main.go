package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/KnoblauchPilze/go-game/cmd/server/routes"
	"github.com/KnoblauchPilze/go-game/pkg/auth"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/KnoblauchPilze/go-game/pkg/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger.Configure(logger.Configuration{
		Service: "server",
	})

	port := 3000
	if len(os.Args) > 1 {
		if maybePort, err := strconv.Atoi(os.Args[1]); err != nil {
			logger.Warnf("Ignoring provided port \"%s\" (err: %v)", os.Args[1], err)
		} else if maybePort <= 0 || maybePort > 65535 {
			logger.Warnf("Ignoring provided port \"%d\" (err: not in range ]0; 65535])", maybePort)
		} else {
			port = maybePort
		}
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	udb := users.NewUserManager()
	tokens := auth.NewAuthenticater()

	r.Mount(routes.SignUpURLRoute, routes.SignUpRouter(udb))
	r.Mount(routes.LoginURLRoute, routes.LoginRouter(udb, tokens))
	r.Mount(routes.UsersURLRoute, routes.UsersRouter(udb, tokens))

	logger.Infof("Starting server on port %d...", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
