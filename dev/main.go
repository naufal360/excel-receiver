package main

import (
	"context"
	"excel-receiver/config"
	"excel-receiver/http/api"
	"excel-receiver/provider"
	"excel-receiver/repository"
	"excel-receiver/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
)

func init() {
	if err := config.LoadConfig("."); err != nil {
		log.Fatal(err)
	}

	provider.InitLogDir()
}

func main() {
	logger := provider.NewLogger()

	conn, err := provider.NewArtemis()
	if err != nil {
		log.Fatal(err)
	}
	logger.Infof(provider.DBLog, "Successfully connected to Artemis")

	db, err := provider.NewMysql()
	if err != nil {
		log.Fatal(err)
	}

	logger.Infof(provider.DBLog, "Successfully connected to DB")

	tokenRepo := repository.NewToken(db, logger)
	requestRepo := repository.NewRequest(db, logger)
	queueRepo := repository.NewQueueArtemis(conn, logger)

	sendRequestService := service.NewSendRequestService(logger, queueRepo, requestRepo)

	app := api.NewApp(logger, sendRequestService, tokenRepo)
	addr := fmt.Sprintf(":%v", config.Configuration.Server.Port)
	server, err := app.CreateServer(addr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		logger.Infof(provider.AppLog, "Server running at: %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf(provider.AppLog, "Server error: %v", err)
		}

	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdownCh
	logger.Infof(provider.AppLog, "Receiving signal: %s", sig)

	ctx, cancel := context.WithTimeout(context.Background(), config.Configuration.Server.ShutdownTimeout)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf(provider.AppLog, "Error during server shutdown: %v", err)
	}

	func(db *sqlx.DB) {
		db.Close()
		logger.Infof(provider.DBLog, "Successfully disconnected from Artemis..")
		logger.Infof(provider.DBLog, "Successfully disconnected from DB..")
		logger.Infof(provider.AppLog, "Successfully shutdown server..")
	}(db)
}
