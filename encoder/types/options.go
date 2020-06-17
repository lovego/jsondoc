package types

type Options struct {
	// quoted causes primitive fields to be encoded inside JSON strings.
	Quoted bool
	// escapeHTML causes '<', '>', and '&' to be escaped in JSON strings.
	EscapeHTML bool
}
