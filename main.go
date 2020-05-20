package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/senseyman/image-media-processor/dto"
	"github.com/senseyman/image-media-processor/server"
	"github.com/sirupsen/logrus"
)

var (
	configPath = "./config.toml"
)

// TODO create service for image processing
// TODO create service for DB
// TODO create service for AWS

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
	logger.Info("Starting application...")

	apiServer := server.NewAPIServer(cfg.Server.ServerPort, logger)

	if err := apiServer.Start(); err != nil {
		logger.Fatalf("Cannot start api server: %v", err)
		panic(err)
	}
}

// config file is required
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

	return logger
}
