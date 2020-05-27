package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/senseyman/image-media-processor/dto/http_request_dto"
	"io"
	"mime/multipart"
	"os"
)

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
