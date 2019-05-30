package notify

import "time"

const (
	Email = iota
	SMS
)

type (
	// Notifier holds the abstract methods that can be performed by the
	// notification engine.
	Notifier interface {
		Send()
		Start()
		Stop()
	}

	Message struct {
		Body       string
		CreatedAt  time.Time
		From       string
		ID         string
		Preference string
		Receiver   string
		Sender     string
		Template   string
		Title      string
		To         string
		Type       int // Browser, SMS, Email etc
	}

	Notify struct {
		ch   chan Message
		quit chan struct{}
	}
)

// Send will deliver a notification
func (n *Notify) Send(msg Message) {
	select {
	case <-n.quit:
	case n.ch <- msg:
	}
}

// New returns a new notify struct
func New() Notify {
	return Notify{}
}

func (n *Notify) Start() {
	go n.loop()
}

func (n *Notify) Stop() {
	close(n.quit)
}

func (n *Notify) loop() {
	for {
		select {
		case <-n.quit:
			return
		case msg, ok := <-n.ch:
			if !ok {
				return
			}
			// Handle different message type.
			switch msg.Type {
			case Email:
			case SMS:
			default:
			}
		}
	}
}
