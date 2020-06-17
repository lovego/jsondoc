package funcs

import (
	"reflect"
	"sync"

	"github.com/lovego/jsondoc/encoder/types"
)

func newStructEncoder(t reflect.Type) encoderFunc {
	se := structEncoder{fields: cachedTypeFields(t)}
	return se.encode
}

type structEncoder struct {
	fields []field
}

func (se structEncoder) encode(buf *types.Buffer, v reflect.Value, opts types.Options) {
	next := byte('{')
FieldLoop:
	for i := range se.fields {
		f := &se.fields[i]

		// Find the nested struct field by following f.index.
		fv := v
		for _, i := range f.index {
			if fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					continue FieldLoop
				}
				fv = fv.Elem()
			}
			fv = fv.Field(i)
		}

		if f.omitEmpty && isEmptyValue(fv) {
			continue
		}
		buf.WriteByte(next)
		next = ','
		if opts.EscapeHTML {
			buf.WriteString(f.nameEscHTML)
		} else {
			buf.WriteString(f.nameNonEsc)
		}
		opts.Quoted = f.quoted
		f.encoder(buf, fv, opts)
	}
	if next == '{' {
		buf.WriteString("{}")
	} else {
		buf.WriteByte('}')
	}
}

var fieldCache sync.Map // map[reflect.Type][]field

// cachedTypeFields is like typeFields but uses a cache to avoid repeated work.
func cachedTypeFields(t reflect.Type) []field {
	if f, ok := fieldCache.Load(t); ok {
		return f.([]field)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.([]field)
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
