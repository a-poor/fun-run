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
		return termenv.ANSIBrightYellow
	case 3:
		return termenv.ANSIBrightRed
	case 4:
		return termenv.ANSIBrightMagenta
	case 5:
		return termenv.ANSIBrightGreen
	default:
		return termenv.ANSIBrightWhite
	}
}

func fmtPrefix(n, t string, w, c int) string {
	// Format the (uncolored) prefix
	s := fmt.Sprintf("%s%s (%s) |", n, strings.Repeat(" ", w-len(n)), t)

	// Create the color-er
	f := termenv.String().Foreground(pickAColor(c))

	// Return the color-ified prefix
	return f.Styled(s)
}

type PrefixWriter struct {
	Prefix []byte
	Writer io.Writer
	sync.Mutex
}

func NewPrefixWriter(name, outType string, nameWidth, color int, write io.Writer) *PrefixWriter {
	return &PrefixWriter{
		Prefix: []byte(fmtPrefix(name, outType, nameWidth, color)),
		Writer: write,
	}
}

func (w *PrefixWriter) Write(p []byte) (int, error) {
	w.Lock()
	defer w.Unlock()
	b := make([]byte, len(w.Prefix)+len(p))
	copy(b, w.Prefix)
	copy(b[len(w.Prefix):], p)

	n, err := w.Writer.Write(b)
	return n - len(w.Prefix), err
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

func (w *SyncWriter) Write(p []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()
	return w.Writer.Write(p)
}
