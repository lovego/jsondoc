package types

import "reflect"

type Options struct {
	// quoted causes primitive fields to be encoded inside JSON strings.
	Quoted bool
	// escapeHTML causes '<', '>', and '&' to be escaped in JSON strings.
	EscapeHTML bool

	// comment to encode inside in struct, slice, array, map values
	comment *string

	// when convert empty slice/map/pointer to non empty ones,
	// record the types has been converted in upper layers to check recursion.
	convertedTypesInUpperLayers []reflect.Type
}

// check recursion when convert empty slice/map/pointer to non empty ones
func (opts *Options) NotRecursion(typ reflect.Type) bool {
	for _, t := range opts.convertedTypesInUpperLayers {
		if t == typ {
			return false
		}
	}
	return true
}

func (opts *Options) AppendConvertedTypeInUpperLayers(typ reflect.Type) {
	opts.convertedTypesInUpperLayers = append(opts.convertedTypesInUpperLayers, typ)
}

// set comment when encode struct field
func (opts *Options) SetComment(comment, commentHTML string) {
	if comment == "" {
		opts.comment = nil
		return
	}
	if opts.EscapeHTML {
		opts.comment = &commentHTML
	} else {
		opts.comment = &comment
	}
}

func (opts *Options) WriteCommentIfPresent(buf *Buffer) {
	if opts.comment != nil && *opts.comment != "" {
		buf.WriteString(*opts.comment)
		*opts.comment = "" // reset parent's comment to empty
		opts.comment = nil
	}
}
