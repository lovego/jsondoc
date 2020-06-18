package funcs

import (
	"encoding"
	"reflect"

	"github.com/lovego/jsondoc/encoder/types"
	"github.com/lovego/jsondoc/scanner"
)

// Marshaler is the interface implemented by types that
// can marshal themselves into valid JSON.
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

// newCondAddrEncoder returns an encoder that checks whether its value
// CanAddr and delegates to canAddrEnc if so, else to elseEnc.
func newCondAddrEncoder(canAddrEnc, elseEnc encoderFunc) encoderFunc {
	enc := condAddrEncoder{canAddrEnc: canAddrEnc, elseEnc: elseEnc}
	return enc.encode
}

type condAddrEncoder struct {
	canAddrEnc, elseEnc encoderFunc
}

func (ce condAddrEncoder) encode(buf *types.Buffer, v reflect.Value, opts types.Options) {
	if v.CanAddr() {
		ce.canAddrEnc(buf, v, opts)
	} else {
		ce.elseEnc(buf, v, opts)
	}
}

func marshalerEncoder(buf *types.Buffer, v reflect.Value, opts types.Options) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		v = reflect.New(v.Type().Elem())
	}
	m, ok := v.Interface().(Marshaler)
	if !ok {
		buf.WriteString("null")
		return
	}
	b, err := m.MarshalJSON()
	if err == nil {
		// copy JSON into types.Buffer, checking validity.
		err = scanner.Compact(&buf.Buffer, b, opts.EscapeHTML)
	}
	if err != nil {
		raiseError(&MarshalerError{v.Type(), err})
	}
}

func addrMarshalerEncoder(buf *types.Buffer, v reflect.Value, _ types.Options) {
	va := v.Addr()
	if va.IsNil() {
		buf.WriteString("null")
		return
	}
	m := va.Interface().(Marshaler)
	b, err := m.MarshalJSON()
	if err == nil {
		// copy JSON into types.Buffer, checking validity.
		err = scanner.Compact(&buf.Buffer, b, true)
	}
	if err != nil {
		raiseError(&MarshalerError{v.Type(), err})
	}
}

func textMarshalerEncoder(buf *types.Buffer, v reflect.Value, opts types.Options) {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		v = reflect.New(v.Type().Elem())
	}
	m := v.Interface().(encoding.TextMarshaler)
	b, err := m.MarshalText()
	if err != nil {
		raiseError(&MarshalerError{v.Type(), err})
	}
	encodeStringBytes(&buf.Buffer, b, opts.EscapeHTML)
}

func addrTextMarshalerEncoder(buf *types.Buffer, v reflect.Value, opts types.Options) {
	va := v.Addr()
	if va.IsNil() {
		buf.WriteString("null")
		return
	}
	m := va.Interface().(encoding.TextMarshaler)
	b, err := m.MarshalText()
	if err != nil {
		raiseError(&MarshalerError{v.Type(), err})
	}
	encodeStringBytes(&buf.Buffer, b, opts.EscapeHTML)
}

// A MarshalerError represents an error from calling a MarshalJSON or MarshalText method.
type MarshalerError struct {
	Type reflect.Type
	Err  error
}

func (e *MarshalerError) Error() string {
	return "json: error calling MarshalJSON for type " + e.Type.String() + ": " + e.Err.Error()
}
