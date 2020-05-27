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

// TODO refactor logs
// Function to handle and process user request for resizing image.
// Handle func call image resize service
// Resized image will send to cloud store and save info to DB store
func (s *ApiServerRequestProcessor) HandleResizeRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	answer := &http_response_dto.ResizeImageResponseDto{}
	jsonEncoder := json.NewEncoder(w)
	s.logger.Info("Got user request")

	r.ParseMultipartForm(100 << 20) // max ~ 100 MB
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		answer.ErrCode = utils.ErrFileNotFoundInRequestCode
		answer.ErrMsg = utils.ErrMsgFileNotFoundInRequest
		jsonEncoder.Encode(answer)
		return
	}

	params := r.FormValue("params")
	if len(params) == 0 {
		s.logger.Error(utils.ErrMsgParamsNotSetInRequest)
		w.WriteHeader(http.StatusBadRequest)
		answer.ErrCode = utils.ErrParamsNotSetInRequestCode
		answer.ErrMsg = utils.ErrMsgParamsNotSetInRequest
		jsonEncoder.Encode(answer)
		return
	}

	rDto := http_request_dto.ResizeImageRequestParamsDto{}
	err = json.Unmarshal([]byte(params), &rDto)
	if err != nil {
		s.logger.Error(utils.ErrMsgCannotParseRequestParams)
		w.WriteHeader(http.StatusBadRequest)
		answer.ErrCode = utils.ErrCannotParseRequestParamsCode
		answer.ErrMsg = utils.ErrMsgCannotParseRequestParams
		jsonEncoder.Encode(answer)
		return
	}

	// validate user request after mapping
	err = s.requestValidator.Validate(rDto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		answer.ErrCode = utils.ErrInvalidRequestParamValuesCode
		answer.ErrMsg = fmt.Sprintf("%s: %v", utils.ErrMsgInvalidRequestParamValues, err)
		s.logger.Errorf(answer.ErrMsg)
		jsonEncoder.Encode(answer)
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
	existEl := s.dbStore.GetPicture(imageId, rDto.Width, rDto.Height)
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
		logEntry.Errorf("RequestId %s -> %s: %v", answer.RequestId, utils.ErrMsgCannotResizeImage, err)
		w.WriteHeader(http.StatusInternalServerError)
		answer.ErrCode = utils.ErrCannotResizeImageCode
		answer.ErrMsg = utils.ErrMsgCannotResizeImage
		jsonEncoder.Encode(answer)
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
		logEntry.Errorf("RequestId %s -> %s: %v", answer.RequestId, utils.ErrMsgUploadImage, err)
		w.WriteHeader(http.StatusInternalServerError)
		answer.ErrCode = utils.ErrUploadImageCode
		answer.ErrMsg = utils.ErrMsgUploadImage
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

	err = s.dbStore.Insert(&dto.DbImageStoreDAO{
		UserId:           rDto.UserId,
		PicId:            imageId,
		OriginalImageUrl: answer.OriginalImagePath,
		ResizedImageUrl:  answer.ResizedImagePath,
		ResizedWidth:     rDto.Width,
		ResizedHeight:    rDto.Height,
	})

	if err != nil {
		logEntry.Errorf("RequestId %s -> %s: %v", answer.RequestId, utils.ErrMsgSaveInfoToDB, err)
		w.WriteHeader(http.StatusInternalServerError)
		answer.ErrCode = utils.ErrSaveInfoToDBCode
		answer.ErrMsg = utils.ErrMsgSaveInfoToDB
		jsonEncoder.Encode(answer)
		return
	}

	jsonEncoder.Encode(answer)
}
