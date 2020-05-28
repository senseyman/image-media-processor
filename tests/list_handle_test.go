package tests

import (
	"encoding/json"
	"fmt"
	"github.com/senseyman/image-media-processor/dto/http_response_dto"
	"github.com/senseyman/image-media-processor/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	Cases:
+	- wrong request type
+	- empty request
+	- without userid
+	- without requestid
	- positive
*/

func TestList_WrongRequestType(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, ApiPathList, nil)
	request.Header.Set("Content-type", "application/json")
	response := httptest.NewRecorder()

	ListRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusMethodNotAllowed, response.Code, "Incorrect response status code on wrong request type")
}

func TestList_InvalidRequest_EmptyRequest(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, ApiPathList, nil)
	request.Header.Set("Content-type", "application/json")
	response := httptest.NewRecorder()

	ListRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code, "Incorrect response status code on wrong request type")
	responseDto := http_response_dto.UserImagesListResponseDto{}
	json.Unmarshal(response.Body.Bytes(), &responseDto)

	assert.Equal(t, utils.ErrInvalidRequestParamValuesCode, responseDto.ErrCode, "Wrong err code")
	assert.Contains(t, responseDto.ErrMsg, utils.ErrMsgInvalidRequestParamValues, "Wrong err message")
	assert.Empty(t, responseDto.UserId, "UserId not empty")
	assert.Empty(t, responseDto.RequestId, "RequestId not empty")
	assert.Nil(t, responseDto.Data, "Wrong err message")
}

func TestList_InvalidRequest_WithoutUserId(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/api/v1/list?request_id=asdsad", nil)
	request.Header.Set("Content-type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()

	ListRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code, "Incorrect response status code on wrong request type")
	responseDto := http_response_dto.UserImagesListResponseDto{}
	json.Unmarshal(response.Body.Bytes(), &responseDto)

	assert.Equal(t, utils.ErrInvalidRequestParamValuesCode, responseDto.ErrCode, "Wrong err code")
	assert.Contains(t, responseDto.ErrMsg, utils.ErrMsgInvalidRequestParamValues, "Wrong err message")
	assert.Empty(t, responseDto.UserId, "UserId not empty")
	assert.Empty(t, responseDto.RequestId, "RequestId not empty")
	assert.Nil(t, responseDto.Data, "Wrong err message")
}

func TestList_InvalidRequest_WithoutRequestId(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/api/v1/list?user_id=affwefwef", nil)
	request.Header.Set("Content-type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()

	ListRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code, "Incorrect response status code on wrong request type")
	responseDto := http_response_dto.UserImagesListResponseDto{}
	json.Unmarshal(response.Body.Bytes(), &responseDto)

	assert.Equal(t, utils.ErrInvalidRequestParamValuesCode, responseDto.ErrCode, "Wrong err code")
	assert.Contains(t, responseDto.ErrMsg, utils.ErrMsgInvalidRequestParamValues, "Wrong err message")
	assert.Empty(t, responseDto.UserId, "UserId not empty")
	assert.Empty(t, responseDto.RequestId, "RequestId not empty")
	assert.Nil(t, responseDto.Data, "Wrong err message")
}

func TestList_Positive(t *testing.T) {
	userId := "affwefwef"
	requestId := "adasdd"
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/list?user_id=%s&request_id=%s", userId, requestId), nil)
	request.Header.Set("Content-type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()

	ListRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code, "Incorrect response status code on wrong request type")
	responseDto := http_response_dto.UserImagesListResponseDto{}
	json.Unmarshal(response.Body.Bytes(), &responseDto)

	assert.Empty(t, responseDto.ErrCode, "Error code not empty")
	assert.Empty(t, responseDto.ErrMsg, "Error message not empty")
	assert.Equal(t, userId, responseDto.UserId, "UserId not equal")
	assert.Equal(t, requestId, responseDto.RequestId, "RequestId not equal")
	assert.NotEmpty(t, responseDto.Data, "Empty response data")
	assert.NotEmpty(t, responseDto.Data[0].PicId, "Wrong pic id")
	assert.NotEmpty(t, responseDto.Data[0].Url, "Wrong url")
	assert.NotEmpty(t, responseDto.Data[0].ResizedImages, "Empty resized array")
	assert.NotEmpty(t, responseDto.Data[0].ResizedImages[0].Url, "Empty resized url")
	assert.NotEmpty(t, responseDto.Data[0].ResizedImages[0].Width, "Empty resized width")
	assert.NotEmpty(t, responseDto.Data[0].ResizedImages[0].Height, "Empty resized height")
}
