package service

import (
	"github.com/senseyman/image-media-processor/dto"
	"io"
	"os"
)

type MediaProcessor interface {
	Resize(buffer io.Reader, name string, width, height int) (*dto.FileInfoDto, error)
}

type CloudStore interface {
	Upload(id uint32, userId string, data []*dto.FileInfoDto) (*dto.CloudResponseDto, error)
	Download(url string, userId string, imageId uint32) (*os.File, error)
}

type DbStore interface {
	Insert(storeDto *dto.DbImageStoreDAO) error
	GetImage(picId uint32, width, height int) *dto.DbImageStoreDAO
	GetImageByImageId(picId uint32) *dto.DbImageStoreDAO
	FindAllPictureByUserId(userId string) []*dto.DbImageStoreDAO
}
