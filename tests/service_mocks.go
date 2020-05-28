package tests

import (
	"bytes"
	"fmt"
	"github.com/senseyman/image-media-processor/dto"
	"io"
	"os"
)

type MediaProcessorMock struct {
	ReturnError bool
}

func (m *MediaProcessorMock) Resize(buffer io.Reader, name string, width, height int) (*dto.FileInfoDto, error) {
	if m.ReturnError {
		return nil, fmt.Errorf("AAAAA")
	}
	bytes.NewBuffer([]byte(""))
	return &dto.FileInfoDto{
		Buffer: bytes.NewBuffer([]byte("")),
		Name:   "name",
		Type:   dto.SourceResized,
	}, nil
}

type CloudStoreMock struct {
}

func (c *CloudStoreMock) Upload(id uint32, userId string, data []*dto.FileInfoDto) (*dto.CloudResponseDto, error) {
	return &dto.CloudResponseDto{
		Data: []*dto.FileCloudStoreDto{
			{
				Id:   id,
				Name: "data[0].Name",
				Type: dto.SourceOriginal,
				Url:  "orig_url",
			},
			{
				Id:   id,
				Name: "data[0].Name",
				Type: dto.SourceResized,
				Url:  "resized_url",
			},
		},
	}, nil
}
func (c *CloudStoreMock) Download(url string, userId string, imageId uint32) (*os.File, error) {
	from, err := os.Open(ImageName)
	if err != nil {
		return nil, err
	}
	defer from.Close()

	to, err := os.OpenFile(fmt.Sprintf("test_%s", ImageName), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(to, from)
	if err != nil {
		return nil, err
	}
	return to, nil
}

type DbStoreMock struct {
}

func (d *DbStoreMock) Insert(storeDto *dto.DbImageStoreDAO) error                    { return nil }
func (d *DbStoreMock) GetImage(picId uint32, width, height int) *dto.DbImageStoreDAO { return nil }
func (d *DbStoreMock) GetImageByImageId(picId uint32) *dto.DbImageStoreDAO {
	return &dto.DbImageStoreDAO{
		UserId:           "asdad",
		PicId:            picId,
		OriginalImageUrl: "orig",
		ResizedImageUrl:  "resized",
		ResizedWidth:     10,
		ResizedHeight:    10,
	}
}

func (d *DbStoreMock) FindAllPictureByUserId(userId string) []*dto.DbImageStoreDAO {
	return []*dto.DbImageStoreDAO{
		{
			UserId:           userId,
			PicId:            1,
			OriginalImageUrl: "orig_url",
			ResizedImageUrl:  "resized_url",
			ResizedWidth:     10,
			ResizedHeight:    10,
		},
	}
}
