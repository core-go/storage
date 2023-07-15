package io

import (
	"path/filepath"
	"strings"
)

func Ext(path string) string {
	ext := filepath.Ext(path)
	if strings.HasPrefix(ext, ":") {
		ext = ext[1:]
	}
	return ext
}
