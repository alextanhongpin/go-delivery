package notify

type notify interface {
	Send() error
}

// Notify is a struct
type Notify struct{}

// Send will deliver a notification
func (notify *Notify) Send() error {
	return nil
}

// New returns a new notify struct
func New() *Notify {
	return &Notify{}
}
