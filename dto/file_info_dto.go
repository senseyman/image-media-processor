package dto

import "io"

type SourceType uint

const (
	SourceOriginal SourceType = iota
	SourceResized
)

type FileInfoDto struct {
	Buffer io.Reader
	Name   string
	Type   SourceType
}

type FileCloudStoreDto struct {
	Id   uint32
	Name string
	Type SourceType
	Url  string
}
