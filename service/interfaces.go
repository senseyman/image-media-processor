package service

import (
	"github.com/senseyman/image-media-processor/dto"
	"github.com/senseyman/image-media-processor/dto/response_dto"
	"io"
)

type MediaProcessor interface {
	Resize(buffer io.Reader, name string, width, height int) (*dto.FileInfoDto, error)
}

type CloudStore interface {
	Upload(id uint32, userId string, data []*dto.FileInfoDto) (*response_dto.CloudResponseDto, error)
}
