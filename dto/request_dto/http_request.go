package request_dto

import (
	"encoding/json"
	"github.com/mholt/binding"
	"mime/multipart"
	"net/http"
)

type BaseRequestDto struct {
	UserId    string `json:"user_id" validate:"regexp=^[[:word:]]"`
	RequestId string `json:"request_id" validate:"regexp=^[[:word:]]"`
}

type SizeRequestDto struct {
	Width  int `json:"width" validate:"min=1"`
	Height int `json:"height" validate:"min=1"`
}

type ResizeImageRequestParamsDto struct {
	BaseRequestDto
	SizeRequestDto
}

type ResizeImageRequestDto struct {
	Params ResizeImageRequestParamsDto `json:"params" `
	File   *multipart.FileHeader       `json:"file"`
}

func (cf *BaseRequestDto) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&cf.UserId:    "user_id",
		&cf.RequestId: "request_id",
	}
}

func (f *ResizeImageRequestDto) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&f.File: "file",
		"a-key": binding.Field{
			Form:   "params",
			Binder: f.Bind,
		},
	}
}

func (f *ResizeImageRequestDto) Bind(fieldName string, strVals []string, ee binding.Errors) binding.Errors {
	if fieldName == "params" && strVals != nil {
		v := strVals[0]
		b := new(ResizeImageRequestParamsDto)
		json.Unmarshal([]byte(v), &b)
		f.Params = *b
	}
	return nil
}

type RequestsHistoryListRequestDto struct {
	BaseRequestDto
}
