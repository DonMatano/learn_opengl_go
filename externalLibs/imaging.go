package externallibs

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/disintegration/imaging"
)

type Imaging struct{}

func NewImagingLib() *Imaging {
	return &Imaging{}
}

func (imag Imaging) LoadImage(filePath string) (*image.RGBA, error) {
	src, err := imaging.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file in path %s, got error %v", filePath, err)
	}
	rgba := image.NewRGBA(src.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), src, image.Point{0, 0}, draw.Src)
	return rgba, nil
}

func (imag Imaging) FlipImage(img *image.RGBA) *image.RGBA {
	return (*image.RGBA)(imaging.FlipV(img))
}
