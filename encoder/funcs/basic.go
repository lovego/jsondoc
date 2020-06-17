package funcs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"

	"github.com/lovego/jsondoc/encoder/types"
)

var (
	float32Encoder = (floatEncoder(32)).encode
	float64Encoder = (floatEncoder(64)).encode
	numberType     = reflect.TypeOf(json.Number(""))
)

func boolEncoder(buf *types.Buffer, v reflect.Value, opts types.Options) {
	if opts.Quoted {
		buf.WriteByte('"')
	}
	if v.Bool() {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
	if opts.Quoted {
		buf.WriteByte('"')
	}
}

func intEncoder(buf *types.Buffer, v reflect.Value, opts types.Options) {
	b := strconv.AppendInt(buf.Scratch[:0], v.Int(), 10)
	if opts.Quoted {
		buf.WriteByte('"')
	}
	buf.Write(b)
	if opts.Quoted {
		buf.WriteByte('"')
	}
}

func uintEncoder(buf *types.Buffer, v reflect.Value, opts types.Options) {
	b := strconv.AppendUint(buf.Scratch[:0], v.Uint(), 10)
	if opts.Quoted {
		buf.WriteByte('"')
	}
	buf.Write(b)
	if opts.Quoted {
		buf.WriteByte('"')
	}
}

func stringEncoder(buf *types.Buffer, v reflect.Value, opts types.Options) {
	if v.Type() == numberType {
		numStr := v.String()
		// In Go1.5 the empty string encodes to "0", while this is not a valid number literal
		// we keep compatibility so check validity after this.
		if numStr == "" {
			numStr = "0" // Number's zero-val
		}
		if !isValidNumber(numStr) {
			raiseError(fmt.Errorf("json: invalid number literal %q", numStr))
		}
		buf.WriteString(numStr)
		return
	}
	if opts.Quoted {
		b := new(bytes.Buffer)
		encodeString(b, v.String(), opts.EscapeHTML)
		encodeString(&buf.Buffer, b.String(), opts.EscapeHTML)
	} else {
		encodeString(&buf.Buffer, v.String(), opts.EscapeHTML)
	}
}

type floatEncoder int // number of bits

func (bits floatEncoder) encode(buf *types.Buffer, v reflect.Value, opts types.Options) {
	f := v.Float()
	if math.IsInf(f, 0) || math.IsNaN(f) {
		raiseError(&UnsupportedValueError{v, strconv.FormatFloat(f, 'g', -1, int(bits))})
	}

	// Convert as if by ES6 number to string conversion.
	// This matches most other JSON generators.
	// See golang.org/issue/6384 and golang.org/issue/14135.
	// Like fmt %g, but the exponent cutoffs are different
	// and exponents themselves are not padded to two digits.
	b := buf.Scratch[:0]
	abs := math.Abs(f)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if bits == 64 && (abs < 1e-6 || abs >= 1e21) || bits == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) {
			fmt = 'e'
		}
	}
	b = strconv.AppendFloat(b, f, fmt, -1, int(bits))
	if fmt == 'e' {
		// clean up e-09 to e-9
		n := len(b)
		if n >= 4 && b[n-4] == 'e' && b[n-3] == '-' && b[n-2] == '0' {
			b[n-2] = b[n-1]
			b = b[:n-1]
		}
	}

	if opts.Quoted {
		buf.WriteByte('"')
	}
	buf.Write(b)
	if opts.Quoted {
		buf.WriteByte('"')
	}
}

type UnsupportedValueError struct {
	Value reflect.Value
	Str   string
}

func (e *UnsupportedValueError) Error() string {
	return "json: unsupported value: " + e.Str
}
