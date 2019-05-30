package omni

import "time"

type Message struct {
	Type          string
	ID            string
	RequestID     string
	NotBeforeTime time.Time
	IssuedAt      time.Time
	From          ContactMechanism
	To            ContactMechanism
	Body          Body
}
