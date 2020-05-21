package server

import (
	"encoding/json"
	"fmt"
	"github.com/senseyman/image-media-processor/dto/request_dto"
	"github.com/senseyman/image-media-processor/dto/response_dto"
	"github.com/senseyman/image-media-processor/utils"
	"net/http"
)

// TODO add mode logging
// Function to handle and process user request_dto for resizing image.
// Handle func call image resize service
// Resized image will send to cloud store and save info to DB store
func (s *APIServer) handleResizeRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	response := &response_dto.ResizeImageResponseDto{}

	s.logger.Info("Got user request")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.ErrCode = utils.ErrInvalidMethodParamsCode
		s.sendResponse(w, response)
		return
	}

	// parse request params
	params := r.Form["params"]
	if len(params) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		response.ErrCode = utils.ErrInvalidMethodParamsCode
		s.sendResponse(w, response)
		return
	}

	requestDto := request_dto.ResizeImageRequestDto{}

	err = json.Unmarshal([]byte(params[0]), &requestDto)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.ErrCode = utils.ErrDecodeMethodParamsCode
		s.sendResponse(w, response)
		return
	}

	if ok, errCode := checkRequest(requestDto); !ok {
		w.WriteHeader(http.StatusBadRequest)
		response.ErrCode = errCode
		s.sendResponse(w, response)
		return
	}

	// save it for response identification on outside
	response.UserId = requestDto.UserId
	response.RequestId = requestDto.RequestId

	// getting file from request form
	file, header, err := r.FormFile("file")
	if err != nil {
		s.logger.Error(fmt.Sprintf("Something went wrongwhile retrieving the file from the form"))
		w.WriteHeader(http.StatusInternalServerError)
		response.ErrCode = utils.ErrSomethingWentWrongCode
		s.sendResponse(w, response)
		return
	}
	defer file.Close()

	filename := header.Filename
	s.logger.Infof("User with id %s send image to resizing. Filename %s", requestDto.UserId, filename)

	// TODO call image processor. Save original and result imgs to cloud. Save img infos to DB

	response.ResizedImagePath = "new path"
	response.OriginalImagePath = "orig path"
	response.ImageId = utils.GenerateImageIdByOriginalName(filename)

	s.sendResponse(w, response)

}

func checkRequest(value request_dto.ResizeImageRequestDto) (bool, int) {
	if value.UserId == "" {
		return false, utils.ErrUserIdIsMissingCode
	}
	if value.RequestId == "" {
		return false, utils.ErrRequestIdIsMissingCode
	}
	if value.Width <= 0 || value.Height <= 0 {
		return false, utils.ErrInvalidResizeParamCode
	}

	return true, 0
}

func (s *APIServer) sendResponse(w http.ResponseWriter, resp *response_dto.ResizeImageResponseDto) {
	v := s.marshalDtoToJson(resp)
	if v == nil {
		s.logger.Error("Cannot send response_dto in marshaled format. Sending plain text")
		fmt.Fprintln(w, fmt.Sprintf("Orig image path %s, new image path %s, errors: code %d and msg %s",
			resp.OriginalImagePath, resp.ResizedImagePath, resp.ErrCode))
	} else {
		w.Write(v)
	}
}
