package server

import (
	"encoding/json"
	"fmt"
	"github.com/mholt/binding"
	"github.com/senseyman/image-media-processor/dto"
	"github.com/senseyman/image-media-processor/dto/request_dto"

	//"github.com/senseyman/image-media-processor/dto/request_dto"
	"github.com/senseyman/image-media-processor/dto/response_dto"
	"github.com/senseyman/image-media-processor/utils"
	"net/http"
)

// TODO refactor logs
// Function to handle and process user request for resizing image.
// Handle func call image resize service
// Resized image will send to cloud store and save info to DB store
func (s *ApiServerRequestProcessor) HandleResizeRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	answer := &response_dto.ResizeImageResponseDto{}
	jsonEncoder := json.NewEncoder(w)
	s.logger.Info("Got user request")

	// mapping http request to dto
	requestDto := new(request_dto.ResizeImageRequestDto)
	if errs := binding.Bind(r, requestDto); errs != nil {
		s.logger.Errorf(fmt.Sprintf("%s: %v", utils.ErrMsgParsingRequestParams, errs))
		w.WriteHeader(http.StatusBadRequest)
		answer.ErrCode = utils.ErrEmptyMethodParamsCode
		answer.ErrMsg = utils.ErrMsgParsingRequestParams
		jsonEncoder.Encode(answer)
		return
	}

	// validate user request after mapping
	err := s.requestValidator.Validate(requestDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		answer.ErrCode = utils.ErrInvalidMethodParamsCode
		answer.ErrMsg = fmt.Sprintf("%s: %v", utils.ErrMsgInvalidUserRequestParams, err)
		s.logger.Errorf(answer.ErrMsg)
		jsonEncoder.Encode(answer)
		return
	}

	// save it for response identification on outside
	answer.UserId = requestDto.Params.UserId
	answer.RequestId = requestDto.Params.RequestId

	if requestDto.File == nil {
		w.WriteHeader(http.StatusBadRequest)
		answer.ErrCode = utils.ErrFileNotFound
		answer.ErrMsg = utils.ErrMsgFileNotFound
		s.logger.Errorf(fmt.Sprintf("UserId %s, requestId %s, err: %s", answer.UserId, answer.RequestId, answer.ErrMsg))
		jsonEncoder.Encode(answer)
		return
	}

	// getting file from request form
	file, err := requestDto.File.Open()
	if err != nil {
		s.logger.Error(fmt.Sprintf("RequestId %s -> %s: %v", answer.RequestId, utils.ErrMsgFileRetrieving, err))
		w.WriteHeader(http.StatusInternalServerError)
		answer.ErrCode = utils.ErrFileRetrievingCode
		answer.ErrMsg = utils.ErrMsgFileRetrieving
		jsonEncoder.Encode(answer)
		return
	}

	filename := requestDto.File.Filename
	s.logger.Infof("User with id %s and requestId %s send image to resizing. Filename %s, width %d, height %d", requestDto.Params.UserId, requestDto.Params.RequestId, filename, requestDto.Params.Width, requestDto.Params.Height)

	resizedFileInfoDto, err := s.imgProcessor.Resize(file, requestDto.File.Filename, requestDto.Params.Width, requestDto.Params.Height)
	// need to close and reopen file for correct using io buffer
	file.Close()

	if err != nil {
		s.logger.Errorf("RequestId %s -> %s: %v", answer.RequestId, utils.ErrMsgResizing, err)
		w.WriteHeader(http.StatusInternalServerError)
		answer.ErrCode = utils.ErrResizingCode
		answer.ErrMsg = utils.ErrMsgResizing
		jsonEncoder.Encode(answer)
		return
	}

	imageId := utils.GenerateImageIdByOriginalName(filename)

	// TODO check in mongo already resized this image with same size params. If exist - send old resized image (get url from DB)

	// reopen file. Skip err if prev were not any error by opening this file
	file, _ = requestDto.File.Open()
	// store images to cloud
	cloudResp, err := s.cloudStore.Upload(imageId, requestDto.Params.UserId, []*dto.FileInfoDto{
		{
			Buffer: file,
			Name:   requestDto.File.Filename,
			Type:   dto.SourceOriginal,
		},
		resizedFileInfoDto,
	})

	if err != nil {
		s.logger.Errorf("RequestId %s -> %s: %v", answer.RequestId, utils.ErrMsgCloudStoring, err)
		w.WriteHeader(http.StatusInternalServerError)
		answer.ErrCode = utils.ErrCloudStoringCode
		answer.ErrMsg = utils.ErrMsgCloudStoring
		jsonEncoder.Encode(answer)
		return
	}

	answer.ImageId = imageId
	for _, c := range cloudResp.Data {
		if c.Type == dto.SourceOriginal {
			answer.OriginalImagePath = c.Url
		} else {
			answer.ResizedImagePath = c.Url
		}
	}

	// TODO store to mongo

	jsonEncoder.Encode(answer)
}
