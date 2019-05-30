package main

const maxQueueLength = 100

type SubscriberServices struct {
	store SubscriberStore
}

func (s *SubscriberService) GetSubscribers() <-chan omni.Notification {
	result := make(<-chan omni.Notification, maxQueueLength)
	go func() {
		ch := s.store.GetSubscribers()
		for subscriber := range ch {
			select {
			case result <- ch:
			}
		}
	}()
	return result
}

func NewSubscriberService(store SubscriberStore) *SubscriberService {
	return &SubscriberService{store}
}

// svc := NewSubscriberService()
// for subscriber := range svc.GetSubscribers() {
//         // publish to message queue.
//         // Mark as delivered to ensure one-time only delivery.
// Mark the last id sent.
// }
