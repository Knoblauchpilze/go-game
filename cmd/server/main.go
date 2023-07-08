package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/KnoblauchPilze/go-game/cmd/server/routes"
	"github.com/KnoblauchPilze/go-game/pkg/db"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/KnoblauchPilze/go-game/pkg/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

const defaultServerPort = 3000

func main() {
	logger.Configure(logger.Configuration{
		Service: "server",
		Level:   logrus.TraceLevel,
	})

	port := getServerPortFromArgs()
	db := createDb()
	repo := users.NewDbRepository(db)
	r := createServerRouter(repo)

	if err := connectToDbAndInstallCleanUp(db); err != nil {
		logger.Fatalf("failed to connect to the db (err: %v)", err)
		return
	}
	defer db.Disconnect()

	logger.Infof("server pid: %d", os.Getpid())
	logger.Infof("starting server on port %d...", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func getServerPortFromArgs() int {
	port := defaultServerPort
	if len(os.Args) > 1 {
		if maybePort, err := strconv.Atoi(os.Args[1]); err != nil {
			logger.Warnf("ignoring provided port \"%s\" (err: %v)", os.Args[1], err)
		} else if maybePort <= 0 || maybePort > 65535 {
			logger.Warnf("ignoring provided port \"%d\" (err: not in range ]0; 65535])", maybePort)
		} else {
			port = maybePort
		}
	}

	return port
}

func createDb() db.Database {
	dbConf := db.NewConfig()
	dbConf.DbHost = "localhost"
	dbConf.DbPort = uint16(5500)
	dbConf.DbName = "user_service_db"
	dbConf.DbUser = "user_service_administrator"
	dbConf.DbPassword = "Ww76hQWbbt7zi2ItM6cNo4YYT"
	dbConf.DbConnectionsPoolSize = 2

	return db.NewPostgresDatabase(dbConf)
}

func createServerRouter(repo users.Repository) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestLogger(routes.TimingLogFormatter{}))
	r.Use(middleware.Recoverer)
	r.Mount("/users", routes.UsersRouter(repo))

	return r
}

func connectToDbAndInstallCleanUp(db db.Database) error {
	if err := db.Connect(); err != nil {
		return err
	}

	// https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-sigint-and-run-a-cleanup-function-i
	interruptChannel := make(chan os.Signal, 2)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interruptChannel

		db.Disconnect()
		os.Exit(1)
	}()

	return nil
}
