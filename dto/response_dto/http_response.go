package response_dto

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
