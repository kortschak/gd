package show

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"

	"github.com/kortschak/gd/internal/enc"
)

// Image renders the given image, title and alt text into the event stream.
func Image(img image.Image, text, title string) error {
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
	return enc.Encode(e)
}

// SVG renders the given SVG image, title and alt text into the event stream.
func SVG(img, text, title string) error {
	e := enc.Event{
		Stream: "image",
		Text:   text,
		Image:  "data:image/svg+xml," + img,
		Title:  title,
	}
	return enc.Encode(e)
}
