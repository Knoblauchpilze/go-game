package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/KnoblauchPilze/go-game/cmd/server/routes"
	"github.com/KnoblauchPilze/go-game/pkg/db"
	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/KnoblauchPilze/go-game/pkg/middleware"
	"github.com/KnoblauchPilze/go-game/pkg/users"
	"github.com/go-chi/chi/v5"
	cmiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const defaultServerPort = 3000

func main() {
	logger.Configure(logger.Configuration{
		Service: "server",
		Level:   logrus.TraceLevel,
	})

	if err := loadConfiguration(); err != nil {
		logger.Errorf("failed to load configuration (err: %v)", err)
		return
	}

	port := viper.GetUint16("Server.Port")
	database := createDb()
	qe := db.NewQueryExecutor(database)
	repo := users.NewDbRepository(qe)
	r := createServerRouter(repo)

	if err := connectToDbAndInstallCleanUp(context.Background(), database); err != nil {
		logger.Errorf("failed to connect to the db (err: %v)", err)
		return
	}
	defer database.Disconnect(context.Background())

	logger.Infof("server pid: %d", os.Getpid())
	logger.Infof("starting server on port %d...", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func loadConfiguration() error {
	// https://github.com/spf13/viper#reading-config-files
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../configs")
	viper.AddConfigPath("../../configs/users")

	// https://github.com/spf13/viper#establishing-defaults
	viper.SetDefault("Server.Port", defaultServerPort)

	viper.SetConfigName("server-dev")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// https://stackoverflow.com/questions/47185318/multiple-config-files-with-go-viper
	viper.SetConfigName("db-dev")
	if err := viper.MergeInConfig(); err != nil {
		return err
	}

	return nil
}

func createDb() db.Database {
	dbConf := db.NewConfig()
	dbConf.DbHost = viper.GetString("Database.Host")
	dbConf.DbPort = viper.GetUint16("Database.Port")
	dbConf.DbName = viper.GetString("Database.Name")
	dbConf.DbUser = viper.GetString("Database.User")
	dbConf.DbPassword = viper.GetString("Database.Password")
	dbConf.DbConnectionsPoolSize = viper.GetUint("Database.ConnectionsPoolSize")

	return db.NewPostgresDatabase(dbConf)
}

func createServerRouter(repo users.Repository) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cmiddleware.Recoverer)
	r.Use(middleware.RequestIdCtx)
	r.Use(middleware.TimingCtx)
	r.Mount("/users", routes.UsersRouter(repo))

	return r
}

func connectToDbAndInstallCleanUp(ctx context.Context, db db.Database) error {
	if err := db.Connect(ctx); err != nil {
		return err
	}

	// https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-sigint-and-run-a-cleanup-function-i
	interruptChannel := make(chan os.Signal, 2)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-interruptChannel

		db.Disconnect(ctx)
		os.Exit(1)
	}()

	return nil
}
