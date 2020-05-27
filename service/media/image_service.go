package media

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/senseyman/image-media-processor/dto"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

// Service for processing images
// Can resize source image to new size
type ImageService struct {
	logger *logrus.Logger
}

func NewImageService(log *logrus.Logger) *ImageService {
	return &ImageService{logger: log}
}

// Function for changing image size (width and height)
// Input params: fileInfo and new size values
// Output - fileInfo and error
// FileInfo include io.Reader and filename
func (i *ImageService) Resize(buffer io.Reader, name string, width, height int) (*dto.FileInfoDto, error) {
	// open file
	src, err := imaging.Decode(buffer)

	if err != nil {
		i.logger.Errorf("failed to open image: %v", err)
		return nil, err
	}

	// call image resizing
	dst := imaging.Resize(src, width, height, imaging.Lanczos)

	format, err := imaging.FormatFromFilename(name)
	if err != nil {
		i.logger.Errorf("failed to get image format: %v", err)
		return nil, err
	}

	buff := new(bytes.Buffer)
	// encode image to buffer
	err = imaging.Encode(buff, dst, format)
	if err != nil {
		fmt.Println("failed to create buffer", err)
	}
	// convert buffer to reader
	reader := bytes.NewReader(buff.Bytes())

	if err != nil {
		i.logger.Errorf("failed to encode dst image: %v", err)
		return nil, err
	}

	origFileExt := strings.Split(name, ".")
	fileNameWithoutExt := strings.ReplaceAll(name, fmt.Sprintf(".%s", origFileExt[1]), "")
	newFileName := fmt.Sprintf("%s_%dx%d.%s", fileNameWithoutExt, width, height, strings.ToLower(format.String()))

	return &dto.FileInfoDto{Buffer: reader, Name: newFileName, Type: dto.SourceResized}, nil

}
