package response_dto

import "github.com/senseyman/image-media-processor/dto"

type CloudResponseDto struct {
	Data []*dto.FileCloudStoreDto
}
