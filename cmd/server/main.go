package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/KnoblauchPilze/go-game/cmd/server/routes"
	"github.com/KnoblauchPilze/go-game/pkg/db"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/KnoblauchPilze/go-game/pkg/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	logger.Configure(logger.Configuration{
		Service: "server",
		Level:   logrus.DebugLevel,
	})

	port := 3000
	if len(os.Args) > 1 {
		if maybePort, err := strconv.Atoi(os.Args[1]); err != nil {
			logger.Warnf("ignoring provided port \"%s\" (err: %v)", os.Args[1], err)
		} else if maybePort <= 0 || maybePort > 65535 {
			logger.Warnf("ignoring provided port \"%d\" (err: not in range ]0; 65535])", maybePort)
		} else {
			port = maybePort
		}
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestLogger(routes.TimingLogFormatter{}))
	r.Use(middleware.Recoverer)

	// repo := users.NewMemoryRepository()
	dbConf := db.NewConfig()
	dbConf.DbHost = "localhost"
	dbConf.DbPort = uint16(5500)
	dbConf.DbName = "user_service_db"
	dbConf.DbUser = "user_service_administrator"
	dbConf.DbPassword = "Ww76hQWbbt7zi2ItM6cNo4YYT"
	dbConf.DbConnectionsPoolSize = 2

	db := db.NewPostgresDatabase(dbConf)
	repo := users.NewDbRepository(db)

	r.Mount("/users", routes.UsersRouter(repo))

	logger.Infof("starting server on port %d...", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
