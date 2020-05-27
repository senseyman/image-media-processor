package utils

// error codes
const (
	ErrUnsupportedMethodCode int = iota + 600
	ErrEmptyRequestCode
	ErrFileNotFoundInRequestCode
	ErrParamsNotSetInRequestCode
	ErrCannotParseRequestParamsCode
	ErrInvalidRequestParamValuesCode
	ErrCannotResizeImageCode
	ErrUploadImageCode
	ErrSaveInfoToDBCode
	ErrCannotGetUserImagesCode
)

// error messages
const (
	ErrMsgUnsupportedMethod         = "Unsupported call method"
	ErrMsgEmptyRequest              = "Empty request"
	ErrMsgFileNotFoundInRequest     = "File not found in request"
	ErrMsgParamsNotSetInRequest     = "Params not set in request"
	ErrMsgCannotParseRequestParams  = "Cannot parse request params"
	ErrMsgInvalidRequestParamValues = "Invalid values in request params"
	ErrMsgCannotResizeImage         = "Cannot resize image"
	ErrMsgUploadImage               = "Cannot upload image to cloud store"
	ErrMsgSaveInfoToDB              = "Cannot save request results to DB"
	ErrMsgCannotGetUserImages       = "Cannot get user images from DB"
)
