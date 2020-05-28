package http_request_dto

type BaseRequestDto struct {
	UserId    string `schema:"user_id" json:"user_id" validate:"regexp=[-a-zA-Z0-9]"`
	RequestId string `schema:"request_id" json:"request_id" validate:"regexp=[-a-zA-Z0-9]"`
}

type SizeRequestDto struct {
	Width  int `json:"width" validate:"min=1"`
	Height int `json:"height" validate:"min=1"`
}

type ResizeImageRequestParamsDto struct {
	BaseRequestDto
	SizeRequestDto
}

type ResizeImageByImageIdRequestParamsDto struct {
	BaseRequestDto
	SizeRequestDto
	ImageId uint32 `schema:"image_id" json:"image_id"`
}

type RequestsHistoryListRequestDto struct {
	BaseRequestDto
}
