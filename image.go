package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"

	"golang.org/x/image/bmp"
	"golang.org/x/image/webp"
)

func init() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
	image.RegisterFormat("bmp", "bmp", bmp.Decode, bmp.DecodeConfig)
	image.RegisterFormat("webp", "webp", webp.Decode, webp.DecodeConfig)
}

type pictable [][][3]uint64

func newPictable(dx int, dy int) pictable {
	pic := make([][][3]uint64, dx) /* type declaration */
	for i := range pic {
		pic[i] = make([][3]uint64, dy) /* again the type? */
		for j := range pic[i] {
			pic[i][j] = [3]uint64{0, 0, 0}
		}
	}
	return pic
}

type imageInfo struct {
	Data     pictable
	Bounds   image.Rectangle
	Filesize uint64
}


func newPictableFromImage(m image.Image) (pictable, error) {
	bounds := m.Bounds()

	avgdata := newPictable(*subdivisions, *subdivisions)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rX := int64(math.Floor((float64(x) / float64(bounds.Max.X)) * float64(*subdivisions)))
			rY := int64(math.Floor((float64(y) / float64(bounds.Max.Y)) * float64(*subdivisions)))

			r, g, b, _ := m.At(x, y).RGBA()
			avgdata[rX][rY][0] += uint64((float32(r) / 65535) * 255)
			avgdata[rX][rY][1] += uint64((float32(g) / 65535) * 255)
			avgdata[rX][rY][2] += uint64((float32(b) / 65535) * 255)
		}
	}

	divisor := uint64((bounds.Max.X / *subdivisions) * (bounds.Max.Y / *subdivisions))
	if divisor == 0 {
		return nil, fmt.Errorf("Image dimensions %d x %d invalid", bounds.Max.X, bounds.Max.Y)
	}

	for rX := 0; rX < *subdivisions; rX++ {
		for rY := 0; rY < *subdivisions; rY++ {
			avgdata[rX][rY][0] = avgdata[rX][rY][0] / divisor
			avgdata[rX][rY][1] = avgdata[rX][rY][1] / divisor
			avgdata[rX][rY][2] = avgdata[rX][rY][2] / divisor
		}
	}

	return avgdata, nil
}
