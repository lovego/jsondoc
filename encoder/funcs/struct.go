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
	buf.WriteByte('{')
	if len(se.fields) > 0 {
		opts.WriteCommentIfPresent(buf)
	}
	needComma := false
	lastFieldOpts := types.Options{}
fieldLoop:
	for i := range se.fields {
		f := &se.fields[i]

		nextLayerOpts := opts // options for next layer

		// Find the nested struct field by following f.index.
		fv := v
		for _, i := range f.index {
			if fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					typ := fv.Type()
					if opts.NotRecursion(typ) {
						fv = reflect.New(typ.Elem())
						nextLayerOpts.AppendConvertedTypeInUpperLayers(typ)
					} else {
						continue fieldLoop
					}
				}
				fv = fv.Elem()
			}
			fv = fv.Field(i)
		}

		if needComma {
			buf.WriteByte(',')
			lastFieldOpts.WriteCommentIfPresent(buf)
		}
		if nextLayerOpts.EscapeHTML {
			buf.WriteString(f.nameEscHTML)
		} else {
			buf.WriteString(f.nameNonEsc)
		}
		nextLayerOpts.Quoted = f.quoted
		nextLayerOpts.SetComment(f.comment, f.commentHTML)

		f.encoder(buf, fv, nextLayerOpts)
		needComma = true
		lastFieldOpts = nextLayerOpts
	}
	lastFieldOpts.WriteCommentIfPresent(buf)
	buf.WriteByte('}')
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
