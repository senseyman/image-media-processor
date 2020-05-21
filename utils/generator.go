package utils

import (
	"hash/fnv"
	"strings"
)

func GenerateImageIdByOriginalName(fileName string) uint32 {
	s := strings.Split(fileName, ".")
	return hash(s[0])
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
