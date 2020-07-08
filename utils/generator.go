package utils

import (
	"hash/fnv"
	"strings"
)

func GenerateImageIdByOriginalName(fileName string) (uint32, error) {
	s := strings.Split(fileName, ".")
	return hash(s[0])
}

func hash(s string) (uint32, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}
