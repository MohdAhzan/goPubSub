package delivery

import (
	"context"
	"log"
	"time"

	"github.com/MohdAhzan/goPubSub/pubsub"
)
type DeliveryOptions struct {
	MaxRetries int
	Backoff    time.Duration
	DLQ        func(pubsub.Message) // Optional
}


func DeliverWithRetry(ctx context.Context, msg pubsub.Message, handler func(pubsub.Message) error, opts DeliveryOptions) {
  

	for attempt := 1; attempt <= opts.MaxRetries; attempt++ {
		err := handler(msg)
		if err == nil {
          ctx.Done()
			return
		}
		log.Printf("Delivery failed (attempt %d): %v", attempt, err)

		select {
		case <-time.After(opts.Backoff * time.Duration(attempt)):
		case <-ctx.Done():
			return
		}
	}

	// Send to DLQ if all retries fail
	if opts.DLQ != nil {
		opts.DLQ(msg)
	}
}


