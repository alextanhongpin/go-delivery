package omni

// Body represents the body of the notification.
type Body struct {
	Subject  string
	Content  []Content // {HTML, "<div>Hello world</div>"} {Text, "hello world"}
	Metadata map[string]interface{}
	Tags     []string
}

func NewBody() *Body {
	return &Body{
		Content:  make([]Content, 0),
		Metadata: make(map[string]interface{}),
	}
}

func (b *Body) AddContent(content Content) {
	b.Content = append(b.Content, content)
}

func (b *Body) AddMetadata(key string, value interface{}) {
	b.Metadata[key] = value
}
