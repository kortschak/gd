package show

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"

	"github.com/kortschak/gd/internal/enc"
)

// Image renders the given image and alt text into the event stream.
func Image(img image.Image, text string) error {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return err
	}
	e := enc.Event{
		Stream: "image",
		Text:   text,
		Image:  "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()),
	}
	return enc.Encode(e)
}
