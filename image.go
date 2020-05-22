package imgdedup

import (
	"errors"
	"fmt"
	"image"
	"math"
	"os"

	// Image format self registers

	// Standard Image Formats
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	// Extended Image Formats
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

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

type ImageInfo struct {
	Data         pictable
	Format       string
	Bounds       image.Rectangle
	Filesize     uint64
	Subdivisions int
}

func NewImageInfo(imgpath string, subdivisions int) (*ImageInfo, error) {
	file, err := os.Open(imgpath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	img, ifmt, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	pict, err := pictableFromImage(img, subdivisions)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	imginfo := &ImageInfo{
		Data:         pict,
		Format:       ifmt,
		Bounds:       img.Bounds(),
		Filesize:     uint64(fi.Size()),
		Subdivisions: subdivisions,
	}

	return imginfo, nil

}

func pictableFromImage(m image.Image, size int) (pictable, error) {
	bounds := m.Bounds()

	avgdata := newPictable(size, size)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rX := int64(math.Floor((float64(x) / float64(bounds.Max.X)) * float64(size)))
			rY := int64(math.Floor((float64(y) / float64(bounds.Max.Y)) * float64(size)))

			r, g, b, _ := m.At(x, y).RGBA()
			avgdata[rX][rY][0] += uint64((float32(r) / 65535) * 255)
			avgdata[rX][rY][1] += uint64((float32(g) / 65535) * 255)
			avgdata[rX][rY][2] += uint64((float32(b) / 65535) * 255)
		}
	}

	divisor := uint64((bounds.Max.X / size) * (bounds.Max.Y / size))
	if divisor == 0 {
		return nil, fmt.Errorf("Image dimensions %d x %d invalid", bounds.Max.X, bounds.Max.Y)
	}

	for rX := 0; rX < size; rX++ {
		for rY := 0; rY < size; rY++ {
			avgdata[rX][rY][0] = avgdata[rX][rY][0] / divisor
			avgdata[rX][rY][1] = avgdata[rX][rY][1] / divisor
			avgdata[rX][rY][2] = avgdata[rX][rY][2] / divisor
		}
	}

	return avgdata, nil
}

// ErrorDissimilarSubdivisions is returned on trying to compare ImageInfo's of different sizes
var ErrorDissimilarSubdivisions = errors.New("diff: cannot compare dissimilar subdivisions")

func Diff(left *ImageInfo, right *ImageInfo) (uint64, error) {
	if left.Subdivisions != right.Subdivisions {
		return 0, ErrorDissimilarSubdivisions
	}

	var xdiff uint64
	for rX := 0; rX < left.Subdivisions; rX++ {
		for rY := 0; rY < left.Subdivisions; rY++ {
			aa := left.Data[rX][rY]
			bb := right.Data[rX][rY]

			xdiff += absdiff(absdiff(absdiff(aa[0], bb[0]), absdiff(aa[1], bb[1])), absdiff(aa[2], bb[2]))
		}
	}
	return xdiff, nil
}

func absdiff(a uint64, b uint64) uint64 {
	return uint64(math.Abs(float64(a) - float64(b)))
}
