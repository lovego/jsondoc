package types

import (
	"bytes"
)

// An Buffer encodes JSON into a bytes.Buffer.
type Buffer struct {
	bytes.Buffer // accumulated output
	Scratch      [64]byte
}
