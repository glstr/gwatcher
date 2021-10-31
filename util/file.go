package util

import (
	"io"
	"os"
)

func OpenReadFile(filePath string) (io.Reader, error) {
	return os.Open(filePath)
}
