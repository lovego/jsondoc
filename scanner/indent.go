// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scanner

import (
	"bytes"
)

// Indent appends to dst an indented form of the JSON-encoded src.
// Each element in a JSON object or array begins on a new,
// indented line beginning with prefix followed by one or more
// copies of indent according to the indentation nesting.
// The data appended to dst does not begin with the prefix nor
// any indentation, to make it easier to embed inside other formatted JSON data.
// Although leading space characters (space, tab, carriage return, newline)
// at the beginning of src are dropped, trailing space characters
// at the end of src are preserved and copied to dst.
// For example, if src has no trailing spaces, neither will dst;
// if src ends in a trailing newline, so will dst.
func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error {
	origLen := dst.Len()
	var scan scanner
	scan.reset()
	depth := 0
	needNewline := false
	for _, c := range src {
		scan.bytes++
		v := scan.step(&scan, c)
		if v == scanSkipSpace {
			continue
		}
		if v == scanError {
			break
		}
		if needNewline {
			needNewline = false
			if v != scanEndObject && v != scanEndArray && v != scanBeginComment {
				newline(dst, prefix, indent, depth)
			}
		}

		switch v {
		// Emit semantically uninteresting bytes
		// (in particular, punctuation in strings) unmodified.
		case scanContinue:
			dst.WriteByte(c)
			continue
		case scanBeginComment:
			dst.WriteString("\t ")
			dst.WriteByte(c)
			continue
		case scanEndComment:
			newline(dst, prefix, indent, depth)
			continue
		}

		// Add spacing around real punctuation.
		switch c {
		case '{', '[':
			dst.WriteByte(c)
			depth++
			// delay newline so that empty object and array are formatted as {} and [].
			needNewline = true

		case ',':
			dst.WriteByte(c)
			// delay newline so that comment are formatted on the same line.
			needNewline = true

		case ':':
			dst.WriteByte(c)
			dst.WriteByte(' ')

		case '}', ']':
			depth--
			newline(dst, prefix, indent, depth)
			dst.WriteByte(c)

		default:
			dst.WriteByte(c)
		}
	}
	if scan.eof() == scanError {
		dst.Truncate(origLen)
		return scan.err
	}
	return nil
}

func newline(dst *bytes.Buffer, prefix, indent string, depth int) {
	dst.WriteByte('\n')
	dst.WriteString(prefix)
	for i := 0; i < depth; i++ {
		dst.WriteString(indent)
	}
}
