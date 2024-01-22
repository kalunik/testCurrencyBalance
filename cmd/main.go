package main

import (
	"context"
	"github.com/kalunik/testCurrencyBalance/config"
	"github.com/kalunik/testCurrencyBalance/internal/app"
	"github.com/kalunik/testCurrencyBalance/pkg/db"
	"github.com/kalunik/testCurrencyBalance/pkg/log"
	"github.com/kalunik/testCurrencyBalance/pkg/nats/consumer"
	"github.com/kalunik/testCurrencyBalance/pkg/nats/producer"
)

func main() {
	logger := log.NewLogger()
	logger.InitLogger()

	configDriver, err := config.LoadNewConfig()
	if err != nil {
		logger.Fatal(err)
	}
	appConfig, err := configDriver.ParseConfig()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("launching app")

	poolConnection, err := db.NewPostgresPoolConnection(context.Background(), appConfig)
	if err != nil {
		logger.Fatal(err)
	}
	defer poolConnection.Close()
	logger.Info("connect to postgres")

	connWriter, err := producer.NewNatsConnection(appConfig, logger)
	if err != nil {
		logger.Fatal(err)
	}
	defer connWriter.Close()
	connReader, err := consumer.NewNatsConnection(appConfig, logger)
	if err != nil {
		logger.Fatal(err)
	}
	defer connReader.Close()

	app.NewApp(poolConnection, connWriter, connReader, logger, appConfig).Run()
}
