package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/senseyman/image-media-processor/dto"
	"github.com/senseyman/image-media-processor/server"
	"github.com/senseyman/image-media-processor/service/db"
	"github.com/senseyman/image-media-processor/service/media"
	"github.com/senseyman/image-media-processor/service/store"
	"github.com/sirupsen/logrus"
)

var (
	configPath = "./config.toml"
)

/*
	TODO create template config file in repo
	TODO add Dockerfile
	TODO add README.md
*/
/*
	main func start application:
	- init config
	- init logger
	- create main services for data processing
	- run api server
*/
func main() {
	cfg := readConfig()
	logger := configureLogger(cfg)

	apiServer := createServer(cfg, logger)

	logger.Info("Starting application...")

	if err := apiServer.Start(); err != nil {
		logger.Fatalf("Cannot start api server: %v", err)
		panic(err)
	}
}

// register all necessary services and return api server instance
func createServer(cfg *dto.Config, logger *logrus.Logger) *server.APIServer {
	logger.Info("Registering services...")
	imgProcessor := media.NewImageService(logger)
	awsService := store.NewAwsService(&cfg.Aws, logger)
	mongoDbService := db.NewMongoDbService(&cfg.MongoDb, logger)
	return server.NewAPIServer(cfg.Server.ServerPort, logger, imgProcessor, awsService, mongoDbService)
}

// Reading configs from config file. File is required
func readConfig() *dto.Config {
	cfg := dto.Config{}
	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		fmt.Printf("Cannot read config file: %v", err)
		panic(err)
	}
	return &cfg
}

// creating logger instance for logging every action in app
func configureLogger(cfg *dto.Config) *logrus.Logger {
	level, err := logrus.ParseLevel(cfg.Server.LogLevel)
	if err != nil {
		fmt.Printf("Cannot configure logger: %v", err)
		panic(err)
	}
	logger := logrus.New()
	logger.SetLevel(level)

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return logger
}
