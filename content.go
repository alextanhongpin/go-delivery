package omni

// Content represents an array of content that can be used to build the final
// content.
type Content struct {
	Type  string
	Value string
}

// NewContent returns a new Content.
func NewContent(typ, value string) Content {
	return Content{Type: typ, Value: value}
}
