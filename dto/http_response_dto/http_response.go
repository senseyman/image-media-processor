package http_response_dto

type BaseResponseDto struct {
	UserId    string `json:"user_id"`
	RequestId string `json:"request_id"`
	ErrCode   int    `json:"err_code"`
	ErrMsg    string `json:"err_msg"`
}

type ResizeImageResponseDto struct {
	BaseResponseDto
	ImageId           uint32 `json:"image_id"`
	OriginalImagePath string `json:"original_image_path"`
	ResizedImagePath  string `json:"resized_image_path"`
}

type UserImagesListResponseDto struct {
	BaseResponseDto
	Data []*UserOriginalImageDbInfoDto `json:"data"`
}

type UserOriginalImageDbInfoDto struct {
	PicId         uint32
	Url           string
	ResizedImages []*UserResizedImageDbInfoDto
}

type UserResizedImageDbInfoDto struct {
	Url    string
	Width  int
	Height int
}
