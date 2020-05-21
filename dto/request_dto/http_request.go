package request_dto

type BaseRequestDto struct {
	UserId    string `json:"user_id"`
	RequestId string `json:"request_id"`
}

type ResizeImageRequestDto struct {
	BaseRequestDto
	Width  int `json:"width"`
	Height int `json:"height"`
}
