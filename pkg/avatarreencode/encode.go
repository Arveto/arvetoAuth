// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package avatarreencode

import (
	"errors"
	"github.com/chai2010/webp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

var (
	ImageNoType      = errors.New("No Content-Type HTTP header")
	ImageTypeUnknown = errors.New("Image unknown media type")
)

// Reencode the avatar of an user.
func Reencode(in io.Reader, t string) ([]byte, error) {
	var decoder func(io.Reader) (image.Image, error)
	switch t {
	case "":
		return nil, ImageNoType
	case "image/png":
		decoder = png.Decode
	case "image/jpeg":
		decoder = jpeg.Decode
	case "image/gif":
		decoder = gif.Decode
	case "image/webp":
		decoder = webp.Decode
	default:
		return nil, ImageTypeUnknown
	}

	img, err := decoder(in)
	if err != nil {
		return nil, err
	}

	return webp.EncodeRGBA(newSquare(img), 100)
}
