package server

import (
	"github.com/senseyman/image-media-processor/dto/http_response_dto"
	"github.com/senseyman/image-media-processor/service"
	"github.com/sirupsen/logrus"
	"gopkg.in/validator.v2"
	"net/http"
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

func writeErrResponseListRequest(w http.ResponseWriter, answer *http_response_dto.UserImagesListResponseDto, serverCode int, errCode int, errMsg string) {
	w.WriteHeader(serverCode)
	answer.ErrCode = errCode
	answer.ErrMsg = errMsg
}

func writeErrResponseResizeRequest(w http.ResponseWriter, answer *http_response_dto.ResizeImageResponseDto, serverCode int, errCode int, errMsg string) {
	w.WriteHeader(serverCode)
	answer.ErrCode = errCode
	answer.ErrMsg = errMsg
}
