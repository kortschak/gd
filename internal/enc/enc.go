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

func Encode(e Event) error {
	mu.Lock()
	defer mu.Unlock()
	_, file, line, _ := runtime.Caller(2)
	e.File = file
	e.Line = line
	pc, _, _, _ := runtime.Caller(1)
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
