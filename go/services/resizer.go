package services

import (
	"image/jpeg"
	"image/png"
	"os"

	"github.com/nfnt/resize"
)

func Resize(originalPath string, newPath string, size uint) (err error) {

	file, err := os.Open(originalPath)
	if err != nil {
		return err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		img, err = png.Decode(file)
		if err != nil {
			return err
		}
	}

	m := resize.Thumbnail(size, size, img, resize.NearestNeighbor)

	out, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, m, nil)
}
