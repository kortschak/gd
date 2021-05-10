// Copyright ©2020 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !gd

package enc

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	mu     sync.Mutex
	events = make(map[int]int)
)

func render(e Event) error {
	i := events[e.Line]
	events[e.Line]++
	switch e.Stream {
	case "stdout":
		fmt.Print(e.Text)
	case "stderr":
		fmt.Fprint(os.Stderr, e.Text)
	case "markdown":
		fmt.Print(e.Text)
	case "image":
		var (
			src    io.Reader
			format string
		)
		switch {
		case strings.HasPrefix(e.Image, "data:image/jpeg;base64,"):
			data := strings.TrimPrefix(e.Image, "data:image/jpeg;base64,")
			src = base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
			format = "jpeg"
		case strings.HasPrefix(e.Image, "data:image/png;base64,"):
			data := strings.TrimPrefix(e.Image, "data:image/png;base64,")
			src = base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
			format = "png"
		case strings.HasPrefix(e.Image, "data:image/svg+xml,"):
			data := strings.TrimPrefix(e.Image, "data:image/svg+xml,")
			src = strings.NewReader(data)
			format = "svg"
		default:
			return fmt.Errorf("unknown image format: %s", e.Image)
		}
		base := filepath.Base(e.File)
		ext := filepath.Ext(base)
		base = base[:len(base)-len(ext)]
		name := fmt.Sprintf("%s_%d_%d.%s", base, e.Line, i, format)

		dst, err := os.Create(name)
		if err != nil {
			return err
		}
		_, err = io.Copy(dst, src)
		return err
	}
	return nil
}
