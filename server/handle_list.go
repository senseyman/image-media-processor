package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"github.com/senseyman/image-media-processor/dto"
	"github.com/senseyman/image-media-processor/dto/http_request_dto"
	"github.com/senseyman/image-media-processor/dto/http_response_dto"
	"github.com/senseyman/image-media-processor/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Function to handle user request for getting list of all his request history
// - requested image and resized results
// - requested resize params
func (s *ApiServerRequestProcessor) HandleListHistoryRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	answer := &http_response_dto.UserImagesListResponseDto{}
	jsonEncoder := json.NewEncoder(w)
	s.logger.Info("Got user request")

	r.ParseForm()
	rDto := http_request_dto.BaseRequestDto{}

	// try to decode request params
	schema.NewDecoder().Decode(&rDto, r.URL.Query())

	// validate user request after mapping
	err := s.requestValidator.Validate(rDto)
	if err != nil {
		errMsg := fmt.Sprintf("%s: %v", utils.ErrMsgInvalidRequestParamValues, err)
		s.logger.Errorf(errMsg)
		writeErrResponseListRequest(w, answer, http.StatusBadRequest, utils.ErrInvalidRequestParamValuesCode, errMsg)
		jsonEncoder.Encode(answer)
		return
	}
	answer.UserId = rDto.UserId
	answer.RequestId = rDto.RequestId

	logEntity := s.logger.WithFields(logrus.Fields{
		"userId":    rDto.UserId,
		"requestId": rDto.RequestId,
	})

	logEntity.Info("Searching user images in DB")
	// try to find all images in DB by userId
	allImgs := s.dbStore.FindAllPictureByUserId(rDto.UserId)

	if allImgs == nil {
		errMsg := fmt.Sprintf("%s: %v", utils.ErrMsgCannotGetUserImages, err)
		logEntity.Errorf(errMsg)
		writeErrResponseListRequest(w, answer, http.StatusInternalServerError, utils.ErrCannotGetUserImagesCode, errMsg)
		jsonEncoder.Encode(answer)
		return
	}

	logEntity.Info("Found records in DB: %d", len(allImgs))
	// collect all user images through imageId
	processDbResponse(answer, allImgs)

	// send answer to caller
	jsonEncoder.Encode(answer)
}

func processDbResponse(resp *http_response_dto.UserImagesListResponseDto, values []*dto.DbImageStoreDAO) {
	for _, img := range values {

		if resp.Data == nil {
			resp.Data = make([]*http_response_dto.UserOriginalImageDbInfoDto, 0)
		}
		exist := findRecById(resp.Data, img.PicId)
		insertValue := &http_response_dto.UserResizedImageDbInfoDto{
			Url:    img.ResizedImageUrl,
			Width:  img.ResizedWidth,
			Height: img.ResizedHeight,
		}
		if exist != nil {
			exist.ResizedImages = append(exist.ResizedImages, insertValue)
		} else {
			resp.Data = append(resp.Data, &http_response_dto.UserOriginalImageDbInfoDto{
				PicId:         img.PicId,
				Url:           img.OriginalImageUrl,
				ResizedImages: []*http_response_dto.UserResizedImageDbInfoDto{insertValue},
			})
		}

	}
}

func findRecById(data []*http_response_dto.UserOriginalImageDbInfoDto, picId uint32) *http_response_dto.UserOriginalImageDbInfoDto {
	for _, d := range data {
		if d.PicId == picId {
			return d
		}
	}
	return nil
}
