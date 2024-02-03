package lib

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func LoadImage(filePath string) (*image.RGBA, error) {
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file in path %s, got error %v", filePath, err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}
	rgb := image.NewRGBA(img.Bounds())
	if rgb.Stride != rgb.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgb, rgb.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgb, nil
}
