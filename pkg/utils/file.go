package utils

import (
	"os"
)

const NormalFilePerm = 0o666

// WriteFile create or overwrite.
func WriteFile(content, path string) error {
	return WriteBytesToFile([]byte(content), path)
}

func WriteBytesToFile(content []byte, path string) error {
	return os.WriteFile(path, content, NormalFilePerm)
}

func ExistFile(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
