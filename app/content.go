package main

type ContentService struct{}

const maxBufferSize = 1000

func main() {
	ch := make(chan Notification, maxBufferSize)
	q := NewQ()
	q.On("go.srv/initiator", func(payload interface{}) {
		// Marshal payload to struct, then send to background channel.
		// ch <- payload.(Notification?)
	})
	svc := NewContentService()
	msgCh := svc.Build(ch)
	for msg := range msgCh {
		// Persist the state.
		// Mark as read.
		// Send to the next process.
	}
}
