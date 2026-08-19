package stdio

import (
	"io"
	"os"
)

// Out returns w, or os.Stdout when w is nil.
func Out(w io.Writer) io.Writer {
	if w == nil {
		return os.Stdout
	}
	return w
}
