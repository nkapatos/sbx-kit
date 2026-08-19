package experimental

import "fmt"

// ErrNotReady is returned by stub entry points until a feature is implemented.
type ErrNotReady struct {
	Feature string
	Track   string
}

func (e ErrNotReady) Error() string {
	return fmt.Sprintf("%s is not implemented yet (track: sbx-kit experimental %s)", e.Feature, e.Track)
}
