// Copyright ©2020 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package enc provides an JSON event stream protocol.
package enc

import (
	"encoding/json"
	"os"
	"path"
	"runtime"
	"sync"
)

var (
	mu  sync.Mutex
	enc *json.Encoder
)

func init() {
	enc = json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
}

func Encode(e Event, depth int) error {
	mu.Lock()
	defer mu.Unlock()
	_, file, line, _ := runtime.Caller(depth + 1)
	e.File = file
	e.Line = line
	pc, _, _, _ := runtime.Caller(depth)
	e.Func = path.Base(runtime.FuncForPC(pc).Name())
	return enc.Encode(e)
}

type Event struct {
	Stream string `json:"stream"`
	File   string `json:"file"`
	Line   int    `json:"line"`
	Func   string `json:"func,omitempty"`
	Text   string `json:"text"`
	Image  string `json:"image,omitempty"`
	Title  string `json:"title,omitempty"`
}
