package http_request_dto

type BaseRequestDto struct {
	UserId    string `schema:"user_id" json:"user_id" validate:"regexp=^[[:word:]]"`
	RequestId string `schema:"request_id" json:"request_id" validate:"regexp=^[[:word:]]"`
}

type SizeRequestDto struct {
	Width  int `json:"width" validate:"min=1"`
	Height int `json:"height" validate:"min=1"`
}

type ResizeImageRequestParamsDto struct {
	BaseRequestDto
	SizeRequestDto
}

type RequestsHistoryListRequestDto struct {
	BaseRequestDto
}
