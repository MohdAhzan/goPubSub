package memory

import (
	"context"
	"log"
	"time"

	"github.com/MohdAhzan/goPubSub/delivery"
	"github.com/MohdAhzan/goPubSub/pubsub"
)


func (m *memoryBackend)startWorker(topic string){
  
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.queues[topic]; exists {
  
    log.Println("WhAT the hechk")
		return // already running
	}

	queue := make(chan pubsub.Message, 100) // buffer size
	ctx, cancel := context.WithCancel(m.ctx)

	m.queues[topic] = queue
	m.workers[topic] = cancel
    log.Println("stajjjjjjjjjj")

	go func() {
		for {
			select {
			case msg := <-queue:
				m.dispatchMessage(ctx, msg)
			case <-ctx.Done():
				return 
			}
		}
	}() 

}

func (m *memoryBackend) dispatchMessage(ctx context.Context, msg pubsub.Message) {
	m.mu.RLock()
	handlers := append([]func(pubsub.Message){}, m.subcribers[msg.Topic]...)

	// user-specific handlers
	if userMap, ok := m.userHandlers[msg.UserID]; ok {
		if userHandlers, ok := userMap[msg.Topic]; ok {
			handlers = append(handlers, userHandlers...)
		}
	}
	m.mu.RUnlock()

	for _, handler := range handlers {
		// Wrap with retry logic
		go delivery.DeliverWithRetry(ctx, msg, func(m pubsub.Message) error {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[panic] recovered in handler: %v", r)
				}
			}()
			handler(m)
			return nil // in future, your handler may return error
		}, delivery.DeliveryOptions{
			MaxRetries: 3,
			Backoff:    500 * time.Millisecond,
			DLQ:        nil, // hook this later
		})
	}
}
