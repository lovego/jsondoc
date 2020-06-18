package funcs

import (
	"encoding/base64"
	"reflect"

	"github.com/lovego/jsondoc/encoder/types"
)

func newSliceEncoder(t reflect.Type) encoderFunc {
	// Byte slices get special treatment; arrays don't.
	if t.Elem().Kind() == reflect.Uint8 {
		p := reflect.PtrTo(t.Elem())
		if !p.Implements(marshalerType) && !p.Implements(textMarshalerType) {
			return encodeByteSlice
		}
	}
	enc := sliceEncoder{newArrayEncoder(t)}
	return enc.encode
}

func newArrayEncoder(t reflect.Type) encoderFunc {
	enc := arrayEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

// sliceEncoder just wraps an arrayEncoder, checking to make sure the value isn't nil.
type sliceEncoder struct {
	arrayEnc encoderFunc
}

func (se sliceEncoder) encode(buf *types.Buffer, v reflect.Value, opts types.Options) {
	if v.Len() == 0 {
		typ := v.Type()
		if opts.NotRecursion(typ) {
			v = reflect.MakeSlice(typ, 1, 1)
			v.Index(0).Set(reflect.Zero(typ.Elem()))
			opts.AppendConvertedTypeInUpperLayers(typ)
		}
	}

	if v.IsNil() {
		buf.WriteString("null")
		return
	}
	se.arrayEnc(buf, v, opts)
}

type arrayEncoder struct {
	elemEnc encoderFunc
}

func (ae arrayEncoder) encode(buf *types.Buffer, v reflect.Value, opts types.Options) {
	buf.WriteByte('[')
	n := v.Len()
	if n > 0 {
		opts.WriteCommentIfPresent(buf)
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			buf.WriteByte(',')
		}
		ae.elemEnc(buf, v.Index(i), opts)
	}
	buf.WriteByte(']')
}

func encodeByteSlice(buf *types.Buffer, v reflect.Value, _ types.Options) {
	if v.IsNil() {
		buf.WriteString("null")
		return
	}
	s := v.Bytes()
	buf.WriteByte('"')
	encodedLen := base64.StdEncoding.EncodedLen(len(s))
	if encodedLen <= len(buf.Scratch) {
		// If the encoded bytes fit in buf.Scratch, avoid an extra
		// allocation and use the cheaper Encoding.Encode.
		dst := buf.Scratch[:encodedLen]
		base64.StdEncoding.Encode(dst, s)
		buf.Write(dst)
	} else if encodedLen <= 1024 {
		// The encoded bytes are short enough to allocate for, and
		// Encoding.Encode is still cheaper.
		dst := make([]byte, encodedLen)
		base64.StdEncoding.Encode(dst, s)
		buf.Write(dst)
	} else {
		// The encoded bytes are too long to cheaply allocate, and
		// Encoding.Encode is no longer noticeably cheaper.
		enc := base64.NewEncoder(base64.StdEncoding, buf)
		enc.Write(s)
		enc.Close()
	}
	buf.WriteByte('"')
}
