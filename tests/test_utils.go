package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/senseyman/image-media-processor/dto/http_request_dto"
	"github.com/senseyman/image-media-processor/server"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

/*

 */
const (
	ApiPathList       = "/api/v1/list"
	ApiPathResize     = "/api/v1/resize"
	ApiPathResizeById = "/api/v1/resize-by-id"

	ImageTag             = "file"
	ImageName            = "image.jpeg"
	UnsupportedImageName = "image.tmp"
)

func ResizeRouter(returnResizeError bool) *mux.Router {
	router := mux.NewRouter()
	logger := logrus.New()
	processor := server.NewApiServerRequestProcessor(logger, &MediaProcessorMock{ReturnError: returnResizeError}, &CloudStoreMock{}, &DbStoreMock{})
	router.HandleFunc(ApiPathResize, processor.HandleResizeRequest).Methods(http.MethodPost)
	return router
}

func ResizeByIdRouterRouter() *mux.Router {
	router := mux.NewRouter()
	logger := logrus.New()
	processor := server.NewApiServerRequestProcessor(logger, &MediaProcessorMock{}, &CloudStoreMock{}, &DbStoreMock{})
	router.HandleFunc(ApiPathResizeById, processor.HandleResizeByIdRequest).Methods(http.MethodPost)
	return router
}

func ListRouter() *mux.Router {
	router := mux.NewRouter()
	logger := logrus.New()
	processor := server.NewApiServerRequestProcessor(logger, &MediaProcessorMock{}, &CloudStoreMock{}, &DbStoreMock{})
	router.HandleFunc(ApiPathList, processor.HandleListHistoryRequest).Methods(http.MethodGet)
	return router
}

func GenerateResizeRequestBody() *http_request_dto.ResizeImageRequestParamsDto {
	return &http_request_dto.ResizeImageRequestParamsDto{
		BaseRequestDto: http_request_dto.BaseRequestDto{
			UserId:    "wsss",
			RequestId: "qqq",
		},
		SizeRequestDto: http_request_dto.SizeRequestDto{
			Width:  10,
			Height: 10,
		},
	}
}

func GenerateResizeByIdRequestBody() *http_request_dto.ResizeImageByImageIdRequestParamsDto {
	return &http_request_dto.ResizeImageByImageIdRequestParamsDto{
		BaseRequestDto: http_request_dto.BaseRequestDto{
			UserId:    "sss",
			RequestId: "ddd",
		},
		SizeRequestDto: http_request_dto.SizeRequestDto{
			Width:  13,
			Height: 13,
		},
		ImageId: 10,
	}
}

func MarshalRequestDto(data interface{}) *bytes.Reader {
	requestByte, _ := json.Marshal(data)

	return bytes.NewReader(requestByte)
}

func prepareRequestValueForResizeApi(requestReader *bytes.Reader, includeImage bool, imageTag, imageName string) (*io.PipeReader, string) {
	body, writer := io.Pipe()

	mwriter := multipart.NewWriter(writer)

	errchan := make(chan error)

	go func() {

		defer close(errchan)
		defer writer.Close()
		defer mwriter.Close()

		if includeImage {
			w, err := mwriter.CreateFormFile(imageTag, imageName)
			if err != nil {
				errchan <- err
				return
			}

			in, err := os.Open(imageName)
			if err != nil {
				errchan <- err
				return
			}
			defer in.Close()

			if written, err := io.Copy(w, in); err != nil {
				errchan <- fmt.Errorf("error copying %s (%d bytes written): %v", imageName, written, err)
				return
			}
		}

		// ***************************************************
		// body
		if requestReader != nil {
			bw, berr := mwriter.CreateFormField("params")
			if berr != nil {
				errchan <- berr
				return
			}
			if written, err := io.Copy(bw, requestReader); err != nil {
				errchan <- fmt.Errorf("error copying requestDto (%d bytes written): %v", written, err)
				return
			}
		}
		// ***************************************************

		if err := mwriter.Close(); err != nil {
			errchan <- err
			return
		}
	}()

	return body, mwriter.FormDataContentType()
}
