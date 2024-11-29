package utils

import (
	"fmt"
	"github.com/h2non/bimg"
	"image"
	"log"
	"os"
	"path/filepath"

	_ "golang.org/x/image/webp"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type ImageMetaData struct {
	Width  int
	Height int
}

func GetImageFilenameWithThumbnailSuffix(filename string) string {
	basename := filepath.Base(filename)
	ext := filepath.Ext(basename)
	return fmt.Sprintf("%s_thumbnail%s", basename[:len(basename)-len(ext)], ext)
}

func GetImageData(file string) (ImageMetaData, error) {
	log.Printf("GetImageType called with file %v\n", file)
	reader, err := os.Open(file)
	if err != nil {
		log.Printf("GetImageType unable to open imageConfig file, error %v\n", err)
		return ImageMetaData{}, err
	}
	defer func(reader *os.File) {
		err := reader.Close()
		if err != nil {
			log.Printf("GetImageType unable to close imageConfig file, error %v\n", err)
		}
	}(reader)

	imageConfig, _, err := image.DecodeConfig(reader)
	if err != nil {
		log.Printf("GetImageType unable to decode imageConfig file, error %v\n", err)
		return ImageMetaData{}, err
	}

	imageData := ImageMetaData{
		Width:  imageConfig.Width,
		Height: imageConfig.Height,
	}

	return imageData, nil
}

func ResizeImage(sourceFilename string, destinationFilename string, width int) error {
	buffer, err := bimg.Read(sourceFilename)
	if err != nil {
		log.Printf("ResizeImage unable to read imageConfig file, error %v\n", err)
		return err
	}

	newImage, err := bimg.NewImage(buffer).Thumbnail(width)
	if err != nil {
		log.Printf("ResizeImage unable to resize %s file, error %v\n", sourceFilename, err)
	}

	err = bimg.Write(destinationFilename, newImage)
	if err != nil {
		return err
	}

	return nil
}
