package controller

import (
	"image"
	"image/color"
	"math"
)

// gaussianBlur applies a gaussian kernel filter on src and return a new blurred image.
func gaussianBlur(src *image.RGBA, ksize float64) *image.RGBA {
	// kernel of gaussian 15x15
	ks := int(ksize)
	k := make([]float64, ks*ks)
	for i := 0; i < ks; i++ {
		for j := 0; j < ks; j++ {
			k[i*ks+j] = math.Exp(-(math.Pow(float64(i)-ksize/2, 2)+math.Pow(float64(j)-ksize/2, 2))/(2*math.Pow(ksize/2, 2))) / 256
		}
	}

	// make an image that is ksize larger than the original
	dst := image.NewRGBA(src.Bounds())

	// apply
	for y := src.Bounds().Min.Y; y < src.Bounds().Max.Y; y++ {
		for x := src.Bounds().Min.X; x < src.Bounds().Max.X; x++ {
			var r, g, b, a float64
			for ky := 0; ky < ks; ky++ {
				for kx := 0; kx < ks; kx++ {
					// get the source pixel
					c := src.At(x+kx-ks/2, y+ky-ks/2)
					r1, g1, b1, a1 := c.RGBA()
					// get the kernel value
					k := k[ky*ks+kx]
					// accumulate
					r += float64(r1) * k
					g += float64(g1) * k
					b += float64(b1) * k
					a += float64(a1) * k
				}
			}
			// set the destination pixel
			dst.Set(x, y, color.RGBA{uint8(r / 273), uint8(g / 273), uint8(b / 273), uint8(a / 273)})
		}
	}
	return dst
}
