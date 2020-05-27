package service

import (
	"github.com/senseyman/image-media-processor/dto"
	"io"
)

type MediaProcessor interface {
	Resize(buffer io.Reader, name string, width, height int) (*dto.FileInfoDto, error)
}

type CloudStore interface {
	Upload(id uint32, userId string, data []*dto.FileInfoDto) (*dto.CloudResponseDto, error)
}

type DbStore interface {
	Insert(storeDto *dto.DbImageStoreDAO) error
	GetImage(picId uint32, width, height int) *dto.DbImageStoreDAO
	FindAllPictureByUserId(userId string) []*dto.DbImageStoreDAO
}
