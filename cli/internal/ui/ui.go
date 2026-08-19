// Package ui is the shared output layer for sbx-kit.
// Commands compute data; this package renders it. A future TUI can consume
// the same result structs without depending on these helpers.
package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

const labelWidth = 12

// Writer is the CLI output channel. All user-facing print goes through here.
type Writer struct {
	Out     io.Writer
	Err     io.Writer
	IsTTY   bool
	NoColor bool
}

// New builds a Writer. Color is off when NO_COLOR is set, TERM is dumb,
// or out is not a terminal.
func New(out, err io.Writer) *Writer {
	noColor := os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb"
	return &Writer{
		Out:     out,
		Err:     err,
		IsTTY:   isTTY(out),
		NoColor: noColor,
	}
}

func isTTY(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// Color reports whether styled output should be used.
func (w *Writer) Color() bool {
	return w != nil && w.IsTTY && !w.NoColor
}

func (w *Writer) out() io.Writer {
	if w == nil || w.Out == nil {
		return os.Stdout
	}
	return w.Out
}

func (w *Writer) err() io.Writer {
	if w == nil || w.Err == nil {
		return os.Stderr
	}
	return w.Err
}

// Header prints a section banner: "==> text".
func (w *Writer) Header(text string) {
	fmt.Fprintf(w.out(), "==> %s\n", text)
}

// Detail prints an indented labeled value.
func (w *Writer) Detail(label, value string) {
	if value == "" {
		return
	}
	fmt.Fprintf(w.out(), "  %-*s  %s\n", labelWidth, label+":", value)
}

// Warn prints a warning to Err.
func (w *Writer) Warn(text string) {
	fmt.Fprintf(w.err(), "warning: %s\n", text)
}

// Empty prints a parenthesized empty-state line, plus an optional hint.
func (w *Writer) Empty(what, hint string) {
	fmt.Fprintf(w.out(), "(no %s)\n", what)
	if hint != "" {
		fmt.Fprintln(w.out(), hint)
	}
}

// Println writes a line to Out.
func (w *Writer) Println(a ...any) {
	fmt.Fprintln(w.out(), a...)
}

// Printf writes formatted text to Out.
func (w *Writer) Printf(format string, a ...any) {
	fmt.Fprintf(w.out(), format, a...)
}

// Table writes aligned columns. headers are uppercased when not already.
func (w *Writer) Table(headers []string, rows [][]string) error {
	tw := tabwriter.NewWriter(w.out(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, strings.Join(headers, "\t"))
	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
	return tw.Flush()
}
