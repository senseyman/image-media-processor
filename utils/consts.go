package utils

// TODO extend codes
// error codes
const (
	ErrInvalidMethodParamsCode = 600
	ErrEmptyMethodParamsCode   = 601
	ErrFileRetrievingCode      = 602
	ErrResizingCode            = 603
	ErrCloudStoringCode        = 604
	ErrFileNotFound            = 605
)

const (
	ErrMsgParsingRequestParams     = "Cannot parse request params"
	ErrMsgInvalidUserRequestParams = "Error while validate user request"
	ErrMsgFileNotFound             = "User request not include file for resizing"
	ErrMsgFileRetrieving           = "Something went wrong while retrieving the file from the form"
	ErrMsgResizing                 = "Error while resizing image"
	ErrMsgCloudStoring             = "Error while storing images to cloud"
)
