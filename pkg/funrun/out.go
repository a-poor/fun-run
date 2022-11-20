package funrun

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/muesli/termenv"
)

func pickAColor(i int) termenv.ANSIColor {
	switch i % 6 {
	case 0:
		return termenv.ANSIBrightBlue
	case 1:
		return termenv.ANSIBrightCyan
	case 2:
		return termenv.ANSIBrightGreen
	case 3:
		return termenv.ANSIBrightRed
	case 4:
		return termenv.ANSIBrightMagenta
	case 5:
		return termenv.ANSIBrightYellow
	default:
		return termenv.ANSIBrightWhite
	}
}

func fmtPrefix(n, t string, w, c int) string {
	// Format the (uncolored) prefix
	s := fmt.Sprintf("%s%s (%s) | ", n, strings.Repeat(" ", w-len(n)), t)

	// Create the color-er
	f := termenv.String().Foreground(pickAColor(c))

	// Return the color-ified prefix
	return f.Styled(s)
}

type PrefixWriter struct {
	Name   string
	Color  termenv.ANSIColor
	Writer io.Writer
	sync.Mutex
}

func NewPrefixWriter(name, outType string, nameWidth, color int, write io.Writer) *PrefixWriter {
	return &PrefixWriter{
		Name:   name,
		Color:  pickAColor(color),
		Writer: write,
	}
}

func (w *PrefixWriter) withColor(s string) string {
	return termenv.
		String().
		Foreground(w.Color).
		Styled(s)
}

func (w *PrefixWriter) Write(p []byte) (int, error) {
	w.Lock()
	defer w.Unlock()

	// Write the prefix
	pfx := w.withColor(w.Name + " | ")
	b := []byte(pfx)
	n, err := w.Writer.Write(b)
	if err != nil {
		return n, err
	}

	// Write the data
	n, err = w.Writer.Write(p)
	if err != nil {
		return n, err
	}

	// Return the number of bytes written
	return n, nil
}

func (w *PrefixWriter) Logln(s string) error {
	p := w.Name + " | "
	c := w.withColor(p + s)
	b := []byte(c)
	_, err := w.Write(b)
	return err
}

func (w *PrefixWriter) Logf(s string, a ...any) error {
	p := w.Name + " | "
	c := w.withColor(p + fmt.Sprintf(s, a...))
	b := []byte(c)
	_, err := w.Writer.Write(b)
	return err
}

func (w *PrefixWriter) Close() error {
	w.Lock()
	defer w.Unlock()
	if c, ok := w.Writer.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

type SyncWriter struct {
	Writer io.Writer
	sync.Mutex
}

func (w *SyncWriter) Write(p []byte) (int, error) {
	w.Lock()
	defer w.Unlock()

	n, err := w.Writer.Write(p)
	if err != nil {
		return n, err
	}

	// Return the number of bytes written
	return n, err
}
