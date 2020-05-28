package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/senseyman/image-media-processor/dto"
	"github.com/senseyman/image-media-processor/dto/http_request_dto"
	"github.com/senseyman/image-media-processor/dto/http_response_dto"
	"github.com/senseyman/image-media-processor/utils"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
		writeErrResponseResizeRequest(w, answer, http.StatusBadRequest, utils.ErrEmptyRequestCode, utils.ErrMsgEmptyRequest)
		jsonEncoder.Encode(answer)
		return
	}
	// getting file from request using tag 'file'
	file, handler, err := r.FormFile("file")
	if err != nil {
		s.logger.Errorf("%s : %v", utils.ErrMsgFileNotFoundInRequest, err)
		writeErrResponseResizeRequest(w, answer, http.StatusBadRequest, utils.ErrFileNotFoundInRequestCode, utils.ErrMsgFileNotFoundInRequest)
		jsonEncoder.Encode(answer)
		return
	}

	// getting params from request using param name 'params'
	params := r.FormValue("params")
	if len(params) == 0 {
		s.logger.Errorf("%s : %v", utils.ErrMsgParamsNotSetInRequest, err)
		writeErrResponseResizeRequest(w, answer, http.StatusBadRequest, utils.ErrParamsNotSetInRequestCode, utils.ErrMsgParamsNotSetInRequest)
		jsonEncoder.Encode(answer)
		return
	}

	// decode params to struct
	rDto := http_request_dto.ResizeImageRequestParamsDto{}
	err = json.Unmarshal([]byte(params), &rDto)
	if err != nil {
		s.logger.Errorf("%s : %v", utils.ErrMsgCannotParseRequestParams, err)
		writeErrResponseResizeRequest(w, answer, http.StatusBadRequest, utils.ErrCannotParseRequestParamsCode, utils.ErrMsgCannotParseRequestParams)
		jsonEncoder.Encode(answer)
		return
	}

	// validate user request after mapping
	err = s.requestValidator.Validate(rDto)
	if err != nil {
		errMsg := fmt.Sprintf("%s: %v", utils.ErrMsgInvalidRequestParamValues, err)
		s.logger.Errorf(errMsg)
		writeErrResponseResizeRequest(w, answer, http.StatusBadRequest, utils.ErrInvalidRequestParamValuesCode, errMsg)
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
	existEl := s.dbStore.GetImage(imageId, rDto.Width, rDto.Height)
	if existEl != nil {
		logEntry.Warn("This picture already processed by the same request params")

		answer.ImageId = imageId
		answer.OriginalImagePath = existEl.OriginalImageUrl
		answer.ResizedImagePath = existEl.ResizedImageUrl
		jsonEncoder.Encode(answer)
		return
	}

	// main workflow
	s.processImageResizeWorkflow(file, handler.Filename, rDto.Width, rDto.Height, imageId, rDto.UserId, w, answer, logEntry, true)

	// send answer to caller
	jsonEncoder.Encode(answer)

}

func (s *ApiServerRequestProcessor) HandleResizeByIdRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	answer := &http_response_dto.ResizeImageResponseDto{}
	jsonEncoder := json.NewEncoder(w)
	s.logger.Info("Got user request")

	rDto := http_request_dto.ResizeImageByImageIdRequestParamsDto{}

	if r.Body == nil {
		s.logger.Error(utils.ErrMsgEmptyRequest)
		writeErrResponseResizeRequest(w, answer, http.StatusBadRequest, utils.ErrEmptyRequestCode, utils.ErrMsgEmptyRequest)
		jsonEncoder.Encode(answer)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&rDto)
	if err != nil {
		s.logger.Errorf("Cannot parse request: %v", err)
		writeErrResponseResizeRequest(w, answer, http.StatusBadRequest, utils.ErrCannotParseRequestParamsCode, utils.ErrMsgCannotParseRequestParams)
		jsonEncoder.Encode(answer)
		return
	}

	// validate user request after mapping
	err = s.requestValidator.Validate(rDto)
	if err != nil {
		errMsg := fmt.Sprintf("%s: %v", utils.ErrMsgInvalidRequestParamValues, err)
		s.logger.Errorf(errMsg)
		writeErrResponseResizeRequest(w, answer, http.StatusBadRequest, utils.ErrInvalidRequestParamValuesCode, errMsg)
		jsonEncoder.Encode(answer)
		return
	}

	// save it for response identification on outside
	answer.UserId = rDto.UserId
	answer.RequestId = rDto.RequestId
	answer.ImageId = rDto.ImageId

	logEntry := s.logger.WithFields(logrus.Fields{
		"UserId":         rDto.UserId,
		"RequestId":      rDto.RequestId,
		"ImageId":        rDto.ImageId,
		"Request width":  rDto.Width,
		"Request height": rDto.Height,
	})

	// check if this image already exist with the same size params
	exist := s.dbStore.GetImage(rDto.ImageId, rDto.Width, rDto.Height)
	if exist != nil {
		logEntry.Warn("Image already processed with this size params")
		answer.OriginalImagePath = exist.OriginalImageUrl
		answer.ResizedImagePath = exist.ResizedImageUrl
		jsonEncoder.Encode(answer)
		return
	}

	// Try to find one image from DB by imageId to get original image url
	img := s.dbStore.GetImageByImageId(rDto.ImageId)
	if img == nil {
		logEntry.Error("This image never processed by user requests")
		writeErrResponseResizeRequest(w, answer, http.StatusBadRequest, utils.ErrImageNotFoundCode, utils.ErrMsgImageNotFound)
		jsonEncoder.Encode(answer)
		return
	}

	// try to download files using image url
	file, err := s.cloudStore.Download(img.OriginalImageUrl, rDto.UserId, rDto.ImageId)
	if err != nil {
		logEntry.Error("Cannot download image from cloud store: %v", err)
		writeErrResponseResizeRequest(w, answer, http.StatusBadRequest, utils.ErrLoadFileCode, utils.ErrMsgLoadFile)
		jsonEncoder.Encode(answer)
		return
	}

	// delete downloaded file from FS
	defer os.Remove(file.Name())

	// main workflow
	s.processImageResizeWorkflow(file, file.Name(), rDto.Width, rDto.Height, rDto.ImageId, rDto.UserId, w, answer, logEntry, false)

	// we don't save original image again to cloud, so need to set to answer original path using info from DB
	answer.OriginalImagePath = img.OriginalImageUrl

	// send answer to caller
	jsonEncoder.Encode(answer)

}

func (s *ApiServerRequestProcessor) resizeImg(
	origFile io.Reader, filename string,
	width, height int,
	w http.ResponseWriter,
	answer *http_response_dto.ResizeImageResponseDto,
	logEntity *logrus.Entry) *dto.FileInfoDto {

	// resizing image with user request params
	resizedFileInfoDto, err := s.imgProcessor.Resize(origFile, filename, width, height)

	if err != nil {
		errMsg := fmt.Sprintf("%s: %v", utils.ErrMsgCannotResizeImage, err)
		logEntity.Errorf(errMsg)
		writeErrResponseResizeRequest(w, answer, http.StatusInternalServerError, utils.ErrCannotResizeImageCode, errMsg)
		return nil
	}

	return resizedFileInfoDto
}

func (s *ApiServerRequestProcessor) uploadFileToCloud(imageId uint32, userId string, upld []*dto.FileInfoDto,
	w http.ResponseWriter,
	answer *http_response_dto.ResizeImageResponseDto,
	logEntity *logrus.Entry) *dto.CloudResponseDto {

	// upload files to cloud
	cloudResp, err := s.cloudStore.Upload(imageId, userId, upld)

	if err != nil {
		logEntity.Errorf("%s: %v", answer.RequestId, utils.ErrMsgUploadImage, err)
		writeErrResponseResizeRequest(w, answer, http.StatusInternalServerError, utils.ErrUploadImageCode, utils.ErrMsgUploadImage)
		return nil
	}

	return cloudResp
}

func (s *ApiServerRequestProcessor) storeToDb(userId string, imageId uint32, origImagePath, resizedImagePath string, width, height int,
	w http.ResponseWriter,
	answer *http_response_dto.ResizeImageResponseDto,
	logEntity *logrus.Entry) error {

	// insert file info to DB
	err := s.dbStore.Insert(&dto.DbImageStoreDAO{
		UserId:           userId,
		PicId:            imageId,
		OriginalImageUrl: origImagePath,
		ResizedImageUrl:  resizedImagePath,
		ResizedWidth:     width,
		ResizedHeight:    height,
	})

	if err != nil {
		logEntity.Errorf("%s: %v", answer.RequestId, utils.ErrMsgSaveInfoToDB, err)
		writeErrResponseResizeRequest(w, answer, http.StatusInternalServerError, utils.ErrSaveInfoToDBCode, utils.ErrMsgSaveInfoToDB)
		return err
	}
	return nil
}

func (s *ApiServerRequestProcessor) processImageResizeWorkflow(
	origFile io.Reader,
	filename string,
	width, height int,
	imageId uint32,
	userId string,
	w http.ResponseWriter,
	answer *http_response_dto.ResizeImageResponseDto,
	logEntity *logrus.Entry,
	saveOriginal bool) {

	// clone reader
	buf, _ := ioutil.ReadAll(origFile)
	bufToUpload := bytes.NewBuffer(buf)
	bufToResize := bytes.NewBuffer(buf)

	// resize image
	resizedImg := s.resizeImg(bufToResize, filename, width, height, w, answer, logEntity)
	if resizedImg == nil {
		return
	}

	var upld []*dto.FileInfoDto
	if saveOriginal {
		upld = []*dto.FileInfoDto{{
			Buffer: bufToUpload,
			Name:   filename,
			Type:   dto.SourceOriginal,
		}, resizedImg}
	} else {
		upld = []*dto.FileInfoDto{resizedImg}
	}

	// call uploading files
	cloudResp := s.uploadFileToCloud(imageId, userId, upld, w, answer, logEntity)

	if cloudResp == nil {
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

	// call storing to DB
	err := s.storeToDb(userId, imageId, answer.OriginalImagePath, answer.ResizedImagePath, width, height, w, answer, logEntity)
	if err != nil {
		return
	}

}
