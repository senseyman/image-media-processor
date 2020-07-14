package tests

import (
	"encoding/json"
	"github.com/senseyman/image-media-processor/dto/http_request_dto"
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
+	- request without image
+	- request without params
+	- request with incorrect params
+		- userId
+		- requestId
+		- width
+		- height
+	- request with unsupported image type
+	- positive
*/

func TestResizeImage_WrongRequestType(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, ApiPathResize, nil)
	request.Header.Set("Content-type", "application/json")
	response := httptest.NewRecorder()

	ResizeRouter(false).ServeHTTP(response, request)

	assert.Equal(t, http.StatusMethodNotAllowed, response.Code, "Incorrect response status code on wrong request type")
}

func TestResizeImage_InvalidParams_EmptyRequest(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, ApiPathResize, nil)
	request.Header.Set("Content-type", "application/json")
	response := httptest.NewRecorder()

	ResizeRouter(false).ServeHTTP(response, request)

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

func TestResizeImage_InvalidParams_WithoutImage(t *testing.T) {
	requestDto := GenerateResizeRequestBody()
	requestReader := MarshalRequestDto(requestDto)
	body, contentType := prepareRequestValueForResizeApi(requestReader, false, ImageTag, ImageName)

	request, _ := http.NewRequest(http.MethodPost, ApiPathResize, body)
	request.Header.Add("Content-Type", contentType)
	response := httptest.NewRecorder()

	ResizeRouter(false).ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code, "Incorrect server response code")
	responseDto := http_response_dto.ResizeImageResponseDto{}
	err := json.Unmarshal(response.Body.Bytes(), &responseDto)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, utils.ErrFileNotFoundInRequestCode, responseDto.ErrCode, "Wrong error code")
	assert.Contains(t, responseDto.ErrMsg, utils.ErrMsgFileNotFoundInRequest, "Wrong error message")
	checkCommonInvalidParamsResponse(t, &responseDto, nil)

}

func TestResizeImage_InvalidParams_WithoutParams(t *testing.T) {
	body, contentType := prepareRequestValueForResizeApi(nil, true, ImageTag, ImageName)

	request, _ := http.NewRequest(http.MethodPost, ApiPathResize, body)
	request.Header.Add("Content-Type", contentType)
	response := httptest.NewRecorder()

	ResizeRouter(false).ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code, "Incorrect server response code")
	responseDto := http_response_dto.ResizeImageResponseDto{}
	err := json.Unmarshal(response.Body.Bytes(), &responseDto)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, utils.ErrParamsNotSetInRequestCode, responseDto.ErrCode, "Wrong error code")
	assert.Contains(t, responseDto.ErrMsg, utils.ErrMsgParamsNotSetInRequest, "Wrong error message")
	checkCommonInvalidParamsResponse(t, &responseDto, nil)

}

func TestResizeImage_InvalidParams_IncorrectUserId(t *testing.T) {
	requestDto := GenerateResizeRequestBody()
	requestDto.UserId = "     ~~~   "
	requestReader := MarshalRequestDto(requestDto)
	body, contentType := prepareRequestValueForResizeApi(requestReader, true, ImageTag, ImageName)

	request, _ := http.NewRequest(http.MethodPost, ApiPathResize, body)
	request.Header.Add("Content-Type", contentType)
	response := httptest.NewRecorder()

	ResizeRouter(false).ServeHTTP(response, request)

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

func TestResizeImage_InvalidParams_IncorrectRequestId(t *testing.T) {
	requestDto := GenerateResizeRequestBody()
	requestDto.RequestId = "     ~~~   "
	requestReader := MarshalRequestDto(requestDto)
	body, contentType := prepareRequestValueForResizeApi(requestReader, true, ImageTag, ImageName)

	request, _ := http.NewRequest(http.MethodPost, ApiPathResize, body)
	request.Header.Add("Content-Type", contentType)
	response := httptest.NewRecorder()

	ResizeRouter(false).ServeHTTP(response, request)

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

func TestResizeImage_InvalidParams_IncorrectWidth(t *testing.T) {
	requestDto := GenerateResizeRequestBody()
	requestDto.Width = 0
	requestReader := MarshalRequestDto(requestDto)
	body, contentType := prepareRequestValueForResizeApi(requestReader, true, ImageTag, ImageName)

	request, _ := http.NewRequest(http.MethodPost, ApiPathResize, body)
	request.Header.Add("Content-Type", contentType)
	response := httptest.NewRecorder()

	ResizeRouter(false).ServeHTTP(response, request)

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

func TestResizeImage_InvalidParams_IncorrectHeight(t *testing.T) {
	requestDto := GenerateResizeRequestBody()
	requestDto.Height = 0
	requestReader := MarshalRequestDto(requestDto)
	body, contentType := prepareRequestValueForResizeApi(requestReader, true, ImageTag, ImageName)

	request, _ := http.NewRequest(http.MethodPost, ApiPathResize, body)
	request.Header.Add("Content-Type", contentType)
	response := httptest.NewRecorder()

	ResizeRouter(false).ServeHTTP(response, request)

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

func TestResizeImage_InvalidParams_UnsupportedImageType(t *testing.T) {
	requestDto := GenerateResizeRequestBody()
	requestReader := MarshalRequestDto(requestDto)
	body, contentType := prepareRequestValueForResizeApi(requestReader, true, ImageTag, UnsupportedImageName)

	request, _ := http.NewRequest(http.MethodPost, ApiPathResize, body)
	request.Header.Add("Content-Type", contentType)
	response := httptest.NewRecorder()

	ResizeRouter(true).ServeHTTP(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code, "Incorrect server response code")
	responseDto := http_response_dto.ResizeImageResponseDto{}
	err := json.Unmarshal(response.Body.Bytes(), &responseDto)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, utils.ErrCannotResizeImageCode, responseDto.ErrCode, "Wrong error code")
	assert.Contains(t, responseDto.ErrMsg, utils.ErrMsgCannotResizeImage, "Wrong error message")
	checkCommonInvalidParamsResponse(t, &responseDto, requestDto)
}

func TestResizeImage_Positive(t *testing.T) {
	requestDto := GenerateResizeRequestBody()
	requestReader := MarshalRequestDto(requestDto)
	body, contentType := prepareRequestValueForResizeApi(requestReader, true, ImageTag, ImageName)

	request, _ := http.NewRequest(http.MethodPost, ApiPathResize, body)
	request.Header.Add("Content-Type", contentType)
	response := httptest.NewRecorder()

	ResizeRouter(false).ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code, "Incorrect server response code")
	responseDto := http_response_dto.ResizeImageResponseDto{}
	err := json.Unmarshal(response.Body.Bytes(), &responseDto)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 0, responseDto.ErrCode, "Error code not zero")
	assert.Empty(t, responseDto.ErrMsg, "Error message not empty")
	assert.Equal(t, requestDto.RequestId, responseDto.RequestId, "RequestId not equal")
	assert.Equal(t, requestDto.UserId, responseDto.UserId, "RequestId not equal")
	assert.NotEmpty(t, responseDto.ResizedImagePath, "ResizedImagePath is empty")
	assert.NotEmpty(t, responseDto.OriginalImagePath, "OriginalImagePath is empty")
	assert.NotEmpty(t, responseDto.ImageId, "ImageId is empty")
}

func checkCommonInvalidParamsResponse(t *testing.T, response *http_response_dto.ResizeImageResponseDto, request *http_request_dto.ResizeImageRequestParamsDto) {
	assert.Empty(t, response.ResizedImagePath, "ResizedImagePath not empty")
	assert.Empty(t, response.OriginalImagePath, "OriginalImagePath not empty")
	assert.Empty(t, response.ImageId, "ImageId not empty")
	if request == nil {
		assert.Empty(t, response.RequestId, "RequestId not empty")
		assert.Empty(t, response.UserId, "UserId not empty")
	} else {
		assert.Equal(t, request.RequestId, response.RequestId, "RequestId not equal")
		assert.Equal(t, request.UserId, response.UserId, "RequestId not equal")
	}
}
