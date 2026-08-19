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

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
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
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
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

func (w *Writer) style(dest io.Writer) (header, label, warn, errStyle lipgloss.Style) {
	plain := lipgloss.NewStyle()
	if w == nil || w.NoColor || !isTTY(dest) {
		return plain, plain, plain, plain
	}
	r := lipgloss.NewRenderer(dest)
	header = lipgloss.NewStyle().Bold(true).Renderer(r)
	label = lipgloss.NewStyle().Faint(true).Renderer(r)
	warn = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true).Renderer(r)
	errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true).Renderer(r)
	return
}

// Header prints a section banner: "==> text".
func (w *Writer) Header(text string) {
	h, _, _, _ := w.style(w.out())
	fmt.Fprintln(w.out(), h.Render("==> "+text))
}

// Detail prints an indented labeled value.
func (w *Writer) Detail(label, value string) {
	if value == "" {
		return
	}
	_, lab, _, _ := w.style(w.out())
	fmt.Fprintf(w.out(), "  %s  %s\n", lab.Render(fmt.Sprintf("%-*s", labelWidth, label+":")), value)
}

// Warn prints a warning to Err.
func (w *Writer) Warn(text string) {
	_, _, warn, _ := w.style(w.err())
	fmt.Fprintln(w.err(), warn.Render("warning: "+text))
}

// ErrorPrefix prints a styled "Error:" prefix then the message to Err.
func (w *Writer) ErrorPrefix(msg string) {
	_, _, _, es := w.style(w.err())
	fmt.Fprintf(w.err(), "%s %s\n", es.Render("Error:"), msg)
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

// Table writes aligned columns.
func (w *Writer) Table(headers []string, rows [][]string) error {
	tw := tabwriter.NewWriter(w.out(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, strings.Join(headers, "\t"))
	for _, row := range rows {
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}
	return tw.Flush()
}
