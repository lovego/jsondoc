// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package json implements encoding and decoding of JSON as defined in
// RFC 7159. The mapping between JSON and Go values is described
// in the documentation for the Marshal and Unmarshal functions.
//
// See "JSON and Go" for an introduction to this package:
// https://golang.org/doc/articles/json_and_go.html
package jsondoc

import (
	"bytes"

	"github.com/lovego/jsondoc/encoder"
	"github.com/lovego/jsondoc/scanner"
)

// MarshalIndent is like json.Marshal but applies Indent to format the output.
// Each JSON element in the output will begin on a new line beginning with prefix
// followed by one or more copies of indent according to the indentation nesting.
func MarshalIndent(v interface{}, escapeHTML bool, prefix, indent string) ([]byte, error) {
	b, err := encoder.Marshal(v, escapeHTML)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = scanner.Indent(&buf, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
