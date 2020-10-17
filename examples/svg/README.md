# SVG example

This is a hello world for rendering SVG images with `gd`. The output
in this README is derived from running [example.go](example.go)
with `gd example.go >README.md`.
```

package main

import (
	"io"
	"log"
	"net/http"
	"strings"

```
The "show" package is a magic package that is resolved by the
`gd` build tool. The full path may also be used, then it would
be "github.com/kortschak/gd/show".
```
	"show"
)

func main() {
	resp, err := http.Get("https://raw.githubusercontent.com/gonum/gonum/master/gopher.svg")
	if err != nil {
		log.Fatal(err)
	}
	var buf strings.Builder
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	show.SVG(buf.String(), "gopher", "Gonum Gopher")
```
> ![gopher](example_35.svg "Gonum Gopher")
```
}
```
