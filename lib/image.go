package lib

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type ImageUtil interface {
	LoadImage(filePath string) (*image.RGBA, error)
	FlipImage(*image.RGBA) *image.RGBA
}

func LoadImage(imageUtil ImageUtil, filePath string) (*image.RGBA, error) {
	rgba, err := imageUtil.LoadImage(filePath)
	// imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file in path %s, got error %v", filePath, err)
	}
	// defer imgFile.Close()
	//
	// img, _, err := image.Decode(imgFile)
	// if err != nil {
	// 	return nil, err
	// }
	return rgba, nil
}

func FlipImage(img *image.RGBA, imageUtil ImageUtil) *image.RGBA {
	rgba := imageUtil.FlipImage(img)
	return rgba
}

func flipPixels(pixels [][]color.Color) [][]color.Color {
	for i := 0; i < len(pixels); i++ {
		pixelRow := pixels[i]
		for j := 0; j < len(pixelRow)/2; j++ {
			k := len(pixelRow) - j - 1
			pixelRow[j], pixelRow[k] = pixelRow[k], pixelRow[j]
		}
	}
	return pixels
}

func convertToPixels(img *image.RGBA) [][]color.Color {
	size := img.Bounds().Size()
	var pixels [][]color.Color
	for i := 0; i < size.X; i++ {
		var colorSlice []color.Color
		for j := 0; j < size.Y; j++ {
			colorSlice = append(colorSlice, img.At(i, j))
		}
		pixels = append(pixels, colorSlice)
	}
	return pixels
}

func convertPixelsToRGBA(pixels [][]color.Color) *image.RGBA {
	rect := image.Rect(0, 0, len(pixels), len(pixels[0]))
	newImage := image.NewRGBA(rect)
	for x := 0; x < len(pixels); x++ {
		for y := 0; y < len(pixels[0]); y++ {
			colorSlice := pixels[x]
			if colorSlice == nil {
				continue
			}
			pixel := pixels[x][y]
			if pixel == nil {
				continue
			}
			original, ok := color.RGBAModel.Convert(pixel).(color.RGBA)
			if ok {
				newImage.Set(x, y, original)
			}
		}
	}
	return newImage
}
