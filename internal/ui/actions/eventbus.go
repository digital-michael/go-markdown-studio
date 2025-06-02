package actions

import "sync"

// EventBus defines a simple pub/sub event system.
type EventBus interface {
	Subscribe(event string, handler func(payload any))
	Unsubscribe(event string, handler func(payload any))
	Publish(event string, payload any)
}

// SimpleEventBus is a thread-safe implementation of EventBus.
type SimpleEventBus struct {
	subscribers map[string][]func(any)
	mu          sync.RWMutex
}

func NewSimpleEventBus() *SimpleEventBus {
	return &SimpleEventBus{
		subscribers: make(map[string][]func(any)),
	}
}

func (eb *SimpleEventBus) Subscribe(event string, handler func(payload any)) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.subscribers[event] = append(eb.subscribers[event], handler)
}

func (eb *SimpleEventBus) Unsubscribe(event string, handler func(payload any)) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	handlers := eb.subscribers[event]
	for i, h := range handlers {
		// Compare function pointers (not always reliable, but works for static handlers)
		if &h == &handler {
			eb.subscribers[event] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

func (eb *SimpleEventBus) Publish(event string, payload any) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()
	for _, handler := range eb.subscribers[event] {
		go handler(payload)
	}
}
