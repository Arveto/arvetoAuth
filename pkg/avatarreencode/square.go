// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package avatarreencode

import (
	"image"
	"image/color"
)

// The size max of an avatar image.
const MaxSize int = 512

type square struct {
	src  image.Image
	rect image.Rectangle
}

func newSquare(img image.Image) image.Image {
	r := img.Bounds()
	if r.Dx() < MaxSize && r.Dy() < MaxSize {
		return img
	}

	return &square{
		src: img,
		rect: r.Intersect(image.Rectangle{
			Min: r.Min,
			Max: image.Point{r.Min.X + MaxSize, r.Min.Y + MaxSize},
		}),
	}
}

func (s *square) At(x, y int) color.Color {
	if (image.Point{x, y}).In(s.rect) {
		return s.src.At(x, y)
	}
	return color.Transparent
}

func (s *square) Bounds() image.Rectangle { return s.rect }
func (s *square) ColorModel() color.Model { return s.src.ColorModel() }
