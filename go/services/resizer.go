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

	useJpeg := true

	img, err := jpeg.Decode(file)
	if err != nil {
		img, err = png.Decode(file)
		if err != nil {
			return err
		}
		useJpeg = false
	}

	m := resize.Thumbnail(size, size, img, resize.NearestNeighbor)

	out, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if useJpeg {
		jpeg.Encode(out, m, nil)
	} else {
		png.Encode(out, m)
	}
	return
}
