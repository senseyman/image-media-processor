package dto

type DbImageStoreDAO struct {
	UserId           string
	PicId            uint32
	OriginalImageUrl string
	ResizedImageUrl  string
	ResizedWidth     int
	ResizedHeight    int
}
