/*{md}
# Images example

This is a small example demonstrating the use of `gd`. The output
in this README is derived from running [example.go](example.go)
with `gd example.go >README.md`.

This comment will be rendered into the Markdown as plain text,
Markdown rendering of a C-style comment requires the `{md}` tag
immediately following the `/*` comment opening.
*/

package main

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"net/http"
	"path"

	/*{md}
	The "show" package is a magic package that is resolved by the
	`gd` build tool. The full path may also be used, then it would
	be "github.com/kortschak/gd/show".
	*/
	"show"

	_ "image/jpeg"
	_ "image/png"
)

var locations = []string{
	"https://blog.golang.org/gopher/glenda.png",
	"https://blog.golang.org/gopher/gopher.png",
	"https://blog.golang.org/gopher/header.jpg",
}

func main() {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, "|URL|Format|File|")
	fmt.Fprintln(&buf, "|---|------|----|")
	for _, url := range locations {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		_, format, err := image.Decode(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(&buf, "|%s|%s|%s|\n", url, format, path.Base(url))
	}
	show.Markdown(buf.String())
}
