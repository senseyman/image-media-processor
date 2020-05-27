package server

import (
	"github.com/senseyman/image-media-processor/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/validator.v2"
)

type ApiServerRequestProcessor struct {
	logger           *logrus.Logger
	requestValidator *validator.Validator
	imgProcessor     service.MediaProcessor
	cloudStore       service.CloudStore
	dbStore          service.DbStore
}

func NewApiServerRequestProcessor(logger *logrus.Logger, imgProcessor service.MediaProcessor, cloudStore service.CloudStore, dbStore service.DbStore) *ApiServerRequestProcessor {
	return &ApiServerRequestProcessor{
		logger:           logger,
		requestValidator: validator.NewValidator(),
		imgProcessor:     imgProcessor,
		cloudStore:       cloudStore,
		dbStore:          dbStore,
	}
}
