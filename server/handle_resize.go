package server

import (
	"encoding/json"
	"fmt"
	"github.com/senseyman/image-media-processor/dto"
	"github.com/senseyman/image-media-processor/dto/http_request_dto"
	"github.com/senseyman/image-media-processor/dto/http_response_dto"
	"github.com/senseyman/image-media-processor/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Function to handle and process user request for resizing image.
// Handle func call image resize service
// Resized image will send to cloud store and save info to DB store
func (s *ApiServerRequestProcessor) HandleResizeRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	answer := &http_response_dto.ResizeImageResponseDto{}
	jsonEncoder := json.NewEncoder(w)
	s.logger.Info("Got user request")

	// max ~ 100 MB
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		s.logger.Errorf("Cannot pars multipart form: %v", err)
		writeErrResponseResizeRequest(w, answer, jsonEncoder, http.StatusBadRequest, utils.ErrEmptyRequestCode, utils.ErrMsgEmptyRequest)
		return
	}
	// getting file from request using tag 'file'
	file, handler, err := r.FormFile("file")
	if err != nil {
		s.logger.Errorf("%s : %v", utils.ErrMsgFileNotFoundInRequest, err)
		writeErrResponseResizeRequest(w, answer, jsonEncoder, http.StatusBadRequest, utils.ErrFileNotFoundInRequestCode, utils.ErrMsgFileNotFoundInRequest)
		return
	}

	// getting params from request using param name 'params'
	params := r.FormValue("params")
	if len(params) == 0 {
		s.logger.Errorf("%s : %v", utils.ErrMsgParamsNotSetInRequest, err)
		writeErrResponseResizeRequest(w, answer, jsonEncoder, http.StatusBadRequest, utils.ErrParamsNotSetInRequestCode, utils.ErrMsgParamsNotSetInRequest)
		return
	}

	// decode params to struct
	rDto := http_request_dto.ResizeImageRequestParamsDto{}
	err = json.Unmarshal([]byte(params), &rDto)
	if err != nil {
		s.logger.Errorf("%s : %v", utils.ErrMsgCannotParseRequestParams, err)
		writeErrResponseResizeRequest(w, answer, jsonEncoder, http.StatusBadRequest, utils.ErrCannotParseRequestParamsCode, utils.ErrMsgCannotParseRequestParams)
		return
	}

	// validate user request after mapping
	err = s.requestValidator.Validate(rDto)
	if err != nil {
		errMsg := fmt.Sprintf("%s: %v", utils.ErrMsgInvalidRequestParamValues, err)
		s.logger.Errorf(errMsg)
		writeErrResponseResizeRequest(w, answer, jsonEncoder, http.StatusBadRequest, utils.ErrInvalidRequestParamValuesCode, errMsg)
		return
	}

	// save it for response identification on outside
	answer.UserId = rDto.UserId
	answer.RequestId = rDto.RequestId

	// generate image id
	imageId := utils.GenerateImageIdByOriginalName(handler.Filename)

	logEntry := s.logger.WithFields(logrus.Fields{
		"UserId":    rDto.UserId,
		"RequestId": rDto.RequestId,
		"Width":     rDto.Width,
		"Height":    rDto.Height,
		"Filename":  handler.Filename,
		"PictureId": imageId,
	})

	logEntry.Info("User send image to resizing")

	// check in DB if this picture already exist with the same resizing params
	// if exist - return known info for this picture
	// else - continue processing request
	existEl := s.dbStore.GetImage(imageId, rDto.Width, rDto.Height)
	if existEl != nil {
		logEntry.Warn("This picture already processed by the same request params")

		answer.ImageId = imageId
		answer.OriginalImagePath = existEl.OriginalImageUrl
		answer.ResizedImagePath = existEl.ResizedImageUrl
		jsonEncoder.Encode(answer)
		return
	}

	// resizing image with user request params
	resizedFileInfoDto, err := s.imgProcessor.Resize(file, handler.Filename, rDto.Width, rDto.Height)
	// need to close and reopen file for correct using io buffer
	file.Close()

	if err != nil {
		logEntry.Errorf("%s: %v", utils.ErrMsgCannotResizeImage, err)
		writeErrResponseResizeRequest(w, answer, jsonEncoder, http.StatusInternalServerError, utils.ErrCannotResizeImageCode, utils.ErrMsgCannotResizeImage)
		return
	}

	// reopen file. Skip err if prev were not any error by opening this file
	file, _ = handler.Open()

	// store images to cloud
	cloudResp, err := s.cloudStore.Upload(imageId, rDto.UserId, []*dto.FileInfoDto{
		{
			Buffer: file,
			Name:   handler.Filename,
			Type:   dto.SourceOriginal,
		},
		resizedFileInfoDto,
	})

	if err != nil {
		logEntry.Errorf("%s: %v", answer.RequestId, utils.ErrMsgUploadImage, err)
		writeErrResponseResizeRequest(w, answer, jsonEncoder, http.StatusInternalServerError, utils.ErrUploadImageCode, utils.ErrMsgUploadImage)
		return
	}

	// prepare http response
	answer.ImageId = imageId
	for _, c := range cloudResp.Data {
		if c.Type == dto.SourceOriginal {
			answer.OriginalImagePath = c.Url
		} else {
			answer.ResizedImagePath = c.Url
		}
	}

	// insert information about processed images to DB
	err = s.dbStore.Insert(&dto.DbImageStoreDAO{
		UserId:           rDto.UserId,
		PicId:            imageId,
		OriginalImageUrl: answer.OriginalImagePath,
		ResizedImageUrl:  answer.ResizedImagePath,
		ResizedWidth:     rDto.Width,
		ResizedHeight:    rDto.Height,
	})

	if err != nil {
		logEntry.Errorf("%s: %v", answer.RequestId, utils.ErrMsgSaveInfoToDB, err)
		writeErrResponseResizeRequest(w, answer, jsonEncoder, http.StatusInternalServerError, utils.ErrSaveInfoToDBCode, utils.ErrMsgSaveInfoToDB)
		return
	}

	// answer to caller
	jsonEncoder.Encode(answer)
}
