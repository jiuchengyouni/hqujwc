package utils

import (
	"crypto/sha1"
	"fmt"
)

func Sha1String(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Implode(arr []string) string {
	result := ""
	for _, s := range arr {
		result += s
	}
	return result
}
