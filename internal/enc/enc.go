package enc

import (
	"encoding/json"
	"os"
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

func Encode(e Event) error {
	mu.Lock()
	defer mu.Unlock()
	_, file, line, _ := runtime.Caller(2)
	e.File = file
	e.Line = line
	return enc.Encode(e)
}

type Event struct {
	Stream string `json:"stream"`
	File   string `json:"file"`
	Line   int    `json:"line"`
	Text   string `json:"text"`
	Image  string `json:"image,omitempty"`
	Title  string `json:"title,omitempty"`
}
