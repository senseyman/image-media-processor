package utils

// error codes
const (
	ErrEmptyRequestCode int = iota + 600
	ErrFileNotFoundInRequestCode
	ErrParamsNotSetInRequestCode
	ErrCannotParseRequestParamsCode
	ErrInvalidRequestParamValuesCode
	ErrCannotResizeImageCode
	ErrUploadImageCode
	ErrSaveInfoToDBCode
	ErrCannotGetUserImagesCode
	ErrImageNotFoundCode
	ErrLoadFileCode
)

// error messages
const (
	ErrMsgEmptyRequest              = "Empty request"
	ErrMsgFileNotFoundInRequest     = "File not found in request"
	ErrMsgParamsNotSetInRequest     = "Params not set in request"
	ErrMsgCannotParseRequestParams  = "Cannot parse request params"
	ErrMsgInvalidRequestParamValues = "Invalid values in request params"
	ErrMsgCannotResizeImage         = "Cannot resize image"
	ErrMsgUploadImage               = "Cannot upload image to cloud store"
	ErrMsgSaveInfoToDB              = "Cannot save request results to DB"
	ErrMsgCannotGetUserImages       = "Cannot get user images from DB"
	ErrMsgImageNotFound             = "Image not found"
	ErrMsgLoadFile                  = "Cannot download file"
)
