// Copyright ©2020 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package show

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/kortschak/gd/internal/enc"
)

// Markdown renders the Markdown text into the event stream.
func Markdown(text string) error {
	e := enc.Event{
		Stream: "markdown",
		Text:   text,
	}
	return enc.Encode(e, 1)
}

// JPEG renders the given image, title and alt text into the event stream as a JPEG.
func JPEG(img image.Image, o *jpeg.Options, text, title string) error {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, o)
	if err != nil {
		return err
	}
	e := enc.Event{
		Stream: "image",
		Text:   text,
		Image:  "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
		Title:  title,
	}
	return enc.Encode(e, 1)
}

// PNG renders the given image, title and alt text into the event stream as a PNG.
func PNG(img image.Image, text, title string) error {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return err
	}
	e := enc.Event{
		Stream: "image",
		Text:   text,
		Image:  "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
		Title:  title,
	}
	return enc.Encode(e, 1)
}

// SVG renders the given SVG image, title and alt text into the event stream.
func SVG(img, text, title string) error {
	e := enc.Event{
		Stream: "image",
		Text:   text,
		Image:  "data:image/svg+xml," + img,
		Title:  title,
	}
	return enc.Encode(e, 1)
}
