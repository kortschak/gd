// Copyright ©2020 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fmt implements formatted I/O to a JSON stream analogous
// the the standard library fmt package. Not all functions make sense
// to use, but all are replicated from the stdlib fmt package.
// Printing functions that print to os.Stdout and os.Stderr are written
// to a JSON stream that is output to os.Stdout. All other output
// io.Writers are treated as normal and all other functions behave as
// the stdlib fmt functions.
package fmt

import (
	"fmt"
	"io"
	"os"

	"github.com/kortschak/gd/internal/enc"
)

// Errorf formats according to a format specifier and returns the string as a
// value that satisfies error.
//
// If the format specifier includes a %w verb with an error operand,
// the returned error will implement an Unwrap method returning the operand. It is
// invalid to include more than one %w verb or to supply it with an operand
// that does not implement the error interface. The %w verb is otherwise
// a synonym for %v.
func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(format, a...)
}

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	switch w {
	case os.Stdout, os.Stderr:
		e := enc.Event{
			Stream: dst(w),
			Text:   fmt.Sprint(a...),
		}
		err = enc.Encode(e, 1)
		return len(e.Text), err
	default:
		return fmt.Fprint(w, a...)
	}
}

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	switch w {
	case os.Stdout, os.Stderr:
		e := enc.Event{
			Stream: dst(w),
			Text:   fmt.Sprintf(format, a...),
		}
		err = enc.Encode(e, 1)
		return len(e.Text), err
	default:
		return fmt.Fprintf(w, format, a...)
	}
}

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	switch w {
	case os.Stdout, os.Stderr:
		e := enc.Event{
			Stream: dst(w),
			Text:   fmt.Sprintln(a...),
		}
		err = enc.Encode(e, 1)
		return len(e.Text), err
	default:
		return fmt.Fprintln(w, a...)
	}
}

// Fscan scans text read from r, storing successive space-separated
// values into successive arguments. Newlines count as space. It
// returns the number of items successfully scanned. If that is less
// than the number of arguments, err will report why.
func Fscan(r io.Reader, a ...interface{}) (n int, err error) {
	return fmt.Fscan(r, a...)
}

// Fscanf scans text read from r, storing successive space-separated
// values into successive arguments as determined by the format. It
// returns the number of items successfully parsed.
// Newlines in the input must match newlines in the format.
func Fscanf(r io.Reader, format string, a ...interface{}) (n int, err error) {
	return fmt.Fscanf(r, format, a...)
}

// Fscanln is similar to Fscan, but stops scanning at a newline and
// after the final item there must be a newline or EOF.
func Fscanln(r io.Reader, a ...interface{}) (n int, err error) {
	return fmt.Fscanln(r, a...)
}

// Print formats using the default formats for its operands and writes to standard output.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func Print(a ...interface{}) (n int, err error) {
	e := enc.Event{
		Stream: dst(os.Stdout),
		Text:   fmt.Sprint(a...),
	}
	err = enc.Encode(e, 1)
	return len(e.Text), err
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
	e := enc.Event{
		Stream: dst(os.Stdout),
		Text:   fmt.Sprintf(format, a...),
	}
	err = enc.Encode(e, 1)
	return len(e.Text), err
}

// Println formats using the default formats for its operands and writes to standard output.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func Println(a ...interface{}) (n int, err error) {
	e := enc.Event{
		Stream: dst(os.Stdout),
		Text:   fmt.Sprintln(a...),
	}
	err = enc.Encode(e, 1)
	return len(e.Text), err
}

// Scan scans text read from standard input, storing successive
// space-separated values into successive arguments. Newlines count
// as space. It returns the number of items successfully scanned.
// If that is less than the number of arguments, err will report why.
func Scan(a ...interface{}) (n int, err error) {
	return fmt.Scan(a...)
}

// Scanf scans text read from standard input, storing successive
// space-separated values into successive arguments as determined by
// the format. It returns the number of items successfully scanned.
// If that is less than the number of arguments, err will report why.
// Newlines in the input must match newlines in the format.
// The one exception: the verb %c always scans the next rune in the
// input, even if it is a space (or tab etc.) or newline.
func Scanf(format string, a ...interface{}) (n int, err error) {
	return fmt.Scanf(format, a...)
}

// Scanln is similar to Scan, but stops scanning at a newline and
// after the final item there must be a newline or EOF.
func Scanln(a ...interface{}) (n int, err error) {
	return fmt.Scanln(a...)
}

// Sprint formats using the default formats for its operands and returns the resulting string.
// Spaces are added between operands when neither is a string.
func Sprint(a ...interface{}) string {
	return fmt.Sprint(a...)
}

// Sprintf formats according to a format specifier and returns the resulting string.
func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

// Sprintln formats using the default formats for its operands and returns the resulting string.
// Spaces are always added between operands and a newline is appended.
func Sprintln(a ...interface{}) string {
	return fmt.Sprintln(a...)
}

// Sscan scans the argument string, storing successive space-separated
// values into successive arguments. Newlines count as space. It
// returns the number of items successfully scanned. If that is less
// than the number of arguments, err will report why.
func Sscan(str string, a ...interface{}) (n int, err error) {
	return fmt.Sscan(str, a...)
}

// Sscanf scans the argument string, storing successive space-separated
// values into successive arguments as determined by the format. It
// returns the number of items successfully parsed.
// Newlines in the input must match newlines in the format.
func Sscanf(str string, format string, a ...interface{}) (n int, err error) {
	return fmt.Sscanf(str, format, a...)
}

// Sscanln is similar to Sscan, but stops scanning at a newline and
// after the final item there must be a newline or EOF.
func Sscanln(str string, a ...interface{}) (n int, err error) {
	return fmt.Sscanln(str, a...)
}

// Formatter is the interface implemented by values with a custom formatter.
// The implementation of Format may call Sprint(f) or Fprint(f) etc.
// to generate its output.
type Formatter = fmt.Formatter

// GoStringer is implemented by any value that has a GoString method,
// which defines the Go syntax for that value.
// The GoString method is used to print values passed as an operand
// to a %#v format.
type GoStringer = fmt.GoStringer

// ScanState represents the scanner state passed to custom scanners.
// Scanners may do rune-at-a-time scanning or ask the ScanState
// to discover the next space-delimited token.
type ScanState = fmt.ScanState

// Scanner is implemented by any value that has a Scan method, which scans
// the input for the representation of a value and stores the result in the
// receiver, which must be a pointer to be useful. The Scan method is called
// for any argument to Scan, Scanf, or Scanln that implements it.
type Scanner = fmt.Scanner

// State represents the printer state passed to custom formatters.
// It provides access to the io.Writer interface plus information about
// the flags and options for the operand's format specifier.
type State = fmt.State

// Stringer is implemented by any value that has a String method,
// which defines the ``native'' format for that value.
// The String method is used to print values passed as an operand
// to any format that accepts a string or to an unformatted printer
// such as Print.
type Stringer = fmt.Stringer

func dst(w io.Writer) string {
	switch w {
	case os.Stdout:
		return "stdout"
	case os.Stderr:
		return "stderr"
	default:
		panic("unexpected writer")
	}
}
