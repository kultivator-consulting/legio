package utils

import (
	"fmt"
	"os"
)

func GetFileStorageName(filename string) string {
	return fmt.Sprintf("file_%s.dat", filename)
}

func GetFileStoragePathname(filename string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("FILE_STORE_PATH"), GetFileStorageName(filename))
}

func GetFileStoragePath(filename string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("FILE_STORE_PATH"), filename)
}
