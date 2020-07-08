package tests

import (
	"encoding/json"
	"github.com/senseyman/image-media-processor/dto/http_response_dto"
	"github.com/senseyman/image-media-processor/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
	Cases
+	- wrong request type
+	- request with empty data
	- request with incorrect params
		- userId
		- requestId
		- width
		- height
	- positive
*/

func TestResizeImageById_WrongRequestType(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, ApiPathResizeById, nil)
	request.Header.Set("Content-type", "application/json")
	response := httptest.NewRecorder()

	ResizeByIdRouterRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusMethodNotAllowed, response.Code, "Incorrect response status code on wrong request type")
}

func TestResizeImageById_InvalidParams_EmptyRequest(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, ApiPathResizeById, nil)
	request.Header.Set("Content-type", "application/json")
	response := httptest.NewRecorder()

	ResizeByIdRouterRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code, "Incorrect server response code")
	responseDto := http_response_dto.ResizeImageResponseDto{}
	err := json.Unmarshal(response.Body.Bytes(), &responseDto)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, utils.ErrEmptyRequestCode, responseDto.ErrCode, "Wrong error code")
	assert.Contains(t, responseDto.ErrMsg, utils.ErrMsgEmptyRequest, "Wrong error message")
	assert.Empty(t, responseDto.RequestId, "RequestId not empty")
	assert.Empty(t, responseDto.UserId, "UserId not empty")
	checkCommonInvalidParamsResponse(t, &responseDto, nil)

}

func TestResizeImageById_InvalidParams_InvalidUserId(t *testing.T) {
	requestDto := GenerateResizeByIdRequestBody()
	requestDto.UserId = "     ~~~   "
	requestReader := MarshalRequestDto(requestDto)

	request, _ := http.NewRequest(http.MethodPost, ApiPathResizeById, requestReader)
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()

	ResizeByIdRouterRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code, "Incorrect server response code")
	responseDto := http_response_dto.ResizeImageResponseDto{}
	err := json.Unmarshal(response.Body.Bytes(), &responseDto)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, utils.ErrInvalidRequestParamValuesCode, responseDto.ErrCode, "Wrong error code")
	assert.Contains(t, responseDto.ErrMsg, utils.ErrMsgInvalidRequestParamValues, "Wrong error message")
	checkCommonInvalidParamsResponse(t, &responseDto, nil)

}

func TestResizeImageById_InvalidParams_InvalidRequestId(t *testing.T) {
	requestDto := GenerateResizeByIdRequestBody()
	requestDto.RequestId = "     ~~~   "
	requestReader := MarshalRequestDto(requestDto)

	request, _ := http.NewRequest(http.MethodPost, ApiPathResizeById, requestReader)
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()

	ResizeByIdRouterRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code, "Incorrect server response code")
	responseDto := http_response_dto.ResizeImageResponseDto{}
	err := json.Unmarshal(response.Body.Bytes(), &responseDto)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, utils.ErrInvalidRequestParamValuesCode, responseDto.ErrCode, "Wrong error code")
	assert.Contains(t, responseDto.ErrMsg, utils.ErrMsgInvalidRequestParamValues, "Wrong error message")
	checkCommonInvalidParamsResponse(t, &responseDto, nil)

}

func TestResizeImageById_InvalidParams_InvalidWidth(t *testing.T) {
	requestDto := GenerateResizeByIdRequestBody()
	requestDto.Width = 0
	requestReader := MarshalRequestDto(requestDto)

	request, _ := http.NewRequest(http.MethodPost, ApiPathResizeById, requestReader)
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()

	ResizeByIdRouterRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code, "Incorrect server response code")
	responseDto := http_response_dto.ResizeImageResponseDto{}
	err := json.Unmarshal(response.Body.Bytes(), &responseDto)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, utils.ErrInvalidRequestParamValuesCode, responseDto.ErrCode, "Wrong error code")
	assert.Contains(t, responseDto.ErrMsg, utils.ErrMsgInvalidRequestParamValues, "Wrong error message")
	checkCommonInvalidParamsResponse(t, &responseDto, nil)

}

func TestResizeImageById_InvalidParams_InvalidHeight(t *testing.T) {
	requestDto := GenerateResizeByIdRequestBody()
	requestDto.Width = 0
	requestReader := MarshalRequestDto(requestDto)

	request, _ := http.NewRequest(http.MethodPost, ApiPathResizeById, requestReader)
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()

	ResizeByIdRouterRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code, "Incorrect server response code")
	responseDto := http_response_dto.ResizeImageResponseDto{}
	err := json.Unmarshal(response.Body.Bytes(), &responseDto)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, utils.ErrInvalidRequestParamValuesCode, responseDto.ErrCode, "Wrong error code")
	assert.Contains(t, responseDto.ErrMsg, utils.ErrMsgInvalidRequestParamValues, "Wrong error message")
	checkCommonInvalidParamsResponse(t, &responseDto, nil)

}

func TestResizeImageById_Positive(t *testing.T) {
	requestDto := GenerateResizeByIdRequestBody()
	requestReader := MarshalRequestDto(requestDto)

	request, _ := http.NewRequest(http.MethodPost, ApiPathResizeById, requestReader)
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()

	ResizeByIdRouterRouter().ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code, "Incorrect server response code")
	//responseDto := http_response_dto.ResizeImageResponseDto{}
	//json.Unmarshal(response.Body.Bytes(), &responseDto)
	//
	//assert.Equal(t, utils.ErrInvalidRequestParamValuesCode, responseDto.ErrCode, "Wrong error code")
	//assert.Contains(t, responseDto.ErrMsg, utils.ErrMsgInvalidRequestParamValues, "Wrong error message")
	//checkCommonInvalidParamsResponse(t, &responseDto, nil)

}
