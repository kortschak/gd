// Copyright ©2020 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gd renders a single Go source file and it's text and graphic
// output into a single Markdown file.
package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kortschak/gd/internal/enc"
)

func main() {
	inline := flag.Bool("inline", false, "render images as inline data: URIs")
	quote := flag.Bool("quote", true, "quote output chunks")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: %[1]s [options] <src.go>\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(2)
	}
	src, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, flag.Arg(0), src, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// Replace the "fmt" and "show" imports with our hooks.
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.IMPORT {
			for _, spec := range gen.Specs {
				imp := spec.(*ast.ImportSpec)
				switch imp.Path.Value {
				case `"fmt"`:
					imp.Path.Value = `"github.com/kortschak/gd/fmt"`
				case `"show"`:
					imp.Path.Value = `"github.com/kortschak/gd/show"`

				}
			}
		}
	}

	// Find C-style comments with leading {md} mark.
	mdText := make(map[int]*ast.Comment)
	for _, c := range f.Comments {
		for _, l := range c.List {
			if strings.HasPrefix(l.Text, "/*{md}\n") {
				mdText[fset.Position(l.Pos()).Line] = l
			}
		}
	}

	events, err := run(fset, f)
	if err != nil {
		log.Fatal(err)
	}

	longTicks := longestTicks(string(src))
	for _, grp := range events {
		for _, e := range grp {
			n := longestTicks(e.Text)
			if n > longTicks {
				longTicks = n
			}
		}
	}
	ticks := strings.Repeat("`", max(longTicks+1, 3))

	rep := strings.NewReplacer("\n", "\n> ")
	sc := bufio.NewScanner(bytes.NewReader(src))
	var line int
	wasComment := true
	for sc.Scan() {
		r, ok := events[line]
		if ok {
			fmt.Println(ticks)
			for i, e := range r {
				if e.File != flag.Arg(0) {
					log.Fatalf("called event generator in dependency file: %s", e.File)
				}
				switch e.Stream {
				case "stdout", "stderr":
					if !strings.HasSuffix(e.Text, "\n") {
						e.Text += "\n"
					}
					if *quote {
						fmt.Printf("> %s%s\n> %s%s\n", ticks, e.Stream, rep.Replace(e.Text), ticks)
					} else {
						fmt.Printf("%s%s\n%s%s\n", ticks, e.Stream, e.Text, ticks)
					}
				case "markdown":
					fmt.Print(e.Text)
				case "image":
					if *inline {
						e.Image = strings.TrimPrefix(e.Image, "data:image/svg+xml,")
						if e.Title == "" {
							fmt.Printf("![%s](%s)\n\n", e.Text, e.Image)
						} else {
							fmt.Printf("![%s](%s %q)\n\n", e.Text, e.Image, e.Title)
						}
					} else {
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
							log.Fatalf("unknown image format: %s", e.Image)
						}
						var name string
						base := filepath.Base(flag.Arg(0))
						ext := filepath.Ext(base)
						base = base[:len(base)-len(ext)]
						if len(r) == 1 {
							name = fmt.Sprintf("%s_%d.%s", base, e.Line, format)
						} else {
							name = fmt.Sprintf("%s_%d_%d.%s", base, e.Line, i, format)
						}

						dst, err := os.Create(name)
						if err != nil {
							log.Fatal(err)
						}
						_, err = io.Copy(dst, src)
						if err != nil {
							log.Fatal(err)
						}
						if *quote {
							fmt.Print("> ")
						}
						if e.Title == "" {
							fmt.Printf("![%s](%s)\n", e.Text, name)
						} else {
							fmt.Printf("![%s](%s %q)\n", e.Text, name, e.Title)
						}
						if len(r) != 1 && i != len(r)-1 {
							fmt.Println()
						}
					}
				}
			}
			fmt.Println(ticks)
		}
		line++
		c, ok := mdText[line]
		if !ok {
			if wasComment {
				fmt.Println(ticks)
			}
			wasComment = false
			fmt.Println(sc.Text())
		} else {
			text := strings.TrimPrefix(c.Text, "/*{md}")
			text = strings.TrimSuffix(text, "*/")
			if !wasComment {
				fmt.Print(ticks)
			} else {
				text = strings.TrimPrefix(text, "\n")
			}
			indent := fset.Position(c.Pos()).Column - 1
			text = strings.Replace(text, "\n"+strings.Repeat("\t", indent), "\n", -1)
			fmt.Printf("%s%s\n", text, ticks)
			n := fset.Position(c.End()).Line - line
			err = skip(n, sc)
			if err != nil {
				log.Fatal(err)
			}
			line += n
			wasComment = false
		}
	}
	if !wasComment {
		fmt.Println(ticks)
	}
}

func longestTicks(s string) int {
	var m, l int
	for _, r := range s {
		if r != '`' {
			if l > m {
				m = l
			}
			l = 0
			continue
		}
		l++
	}
	return m
}

func skip(n int, sc *bufio.Scanner) error {
	for n--; sc.Scan() && n > 0; n-- {
	}
	if n > 0 {
		if sc.Err() != nil {
			return sc.Err()
		}
		return io.ErrUnexpectedEOF
	}
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// run runs the source described by fset and f and collects output events.
func run(fset *token.FileSet, f *ast.File) (map[int][]enc.Event, error) {
	// Retain line numbering to be consistent with the
	// source as given.
	cfg := printer.Config{
		Mode:     printer.UseSpaces | printer.TabIndent | printer.SourcePos,
		Tabwidth: 8,
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	tmp, err := ioutil.TempFile(wd, "gd-*.go")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())
	err = cfg.Fprint(tmp, fset, f)
	if err != nil {
		return nil, err
	}
	err = tmp.Close()
	if err != nil {
		return nil, err
	}

	gorun := exec.Command("go", "run", tmp.Name())
	var buf bytes.Buffer
	gorun.Stdout = &buf
	gorun.Stderr = os.Stderr
	err = gorun.Run()
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(&buf)

	events := make(map[int][]enc.Event)
	for {
		var e enc.Event
		err = dec.Decode(&e)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		events[e.Line] = append(events[e.Line], e)
	}
	return events, nil
}
