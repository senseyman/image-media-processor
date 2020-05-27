package mock

import "github.com/senseyman/image-media-processor/dto"

type CloudStoreMock struct {
}

func (c *CloudStoreMock) Upload(id uint32, userId string, data []*dto.FileInfoDto) (*dto.CloudResponseDto, error) {
	return &dto.CloudResponseDto{
		Data: []*dto.FileCloudStoreDto{
			{
				Id:   id,
				Name: data[0].Name,
				Type: data[0].Type,
				Url:  "orig_url",
			},
			{
				Id:   id,
				Name: data[1].Name,
				Type: data[1].Type,
				Url:  "resized_url",
			},
		},
	}, nil
}

type DbStoreMock struct {
}

func (d *DbStoreMock) Insert(storeDto *dto.DbImageStoreDAO) error                    { return nil }
func (d *DbStoreMock) GetImage(picId uint32, width, height int) *dto.DbImageStoreDAO { return nil }
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
