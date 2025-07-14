package memory

import (
	"sync"

	"github.com/MohdAhzan/goPubSub/pubsub"
)


type memoryBackend struct{

  mu sync.RWMutex
  subcribers  map[string][]func(pubsub.Message)
  userHandlers map[string]map[string][]func(pubsub.Message)

}

func New()pubsub.Backend{
  
  return &memoryBackend{
    subcribers: make(map[string][]func(pubsub.Message)),
    userHandlers: make(map[string]map[string][]func(pubsub.Message)),
  }

}

// type Backend interface {
// 	Publish(ctx context.Context, msg Message) error
// 	Subscribe(ctx context.Context, topic string, handler func(Message)) error
// 	SubscribeUser(ctx context.Context, topic string, handler func(Message)) error
// 	UnSubscribe(ctx context.Context, topic string) error
// 	Close() error
// }

func (m *memoryBackend)
