package memory

import (
	"context"
	"log"
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

func (m *memoryBackend)Publish(ctx context.Context,msg pubsub.Message)error{

  m.mu.RLock()
  defer m.mu.RUnlock()

  //publishing subscribed topic handlers 
  if msg.UserID != ""{
    if userMap,ok:= m.userHandlers[msg.UserID];ok{
      for _,handler:=range userMap[msg.Topic]{
        go handler(msg)
      }
    }
    return nil
  }

  //publishing subscribed topic handlers 
  for _,handler :=range m.subcribers[msg.Topic]{
    go handler(msg)
  }

  return nil
}

func (m *memoryBackend)Subscribe(ctx context.Context,topic string, handler func(pubsub.Message))error{
  m.mu.RLock() 
  defer m.mu.RUnlock()

  go func() {
    <- ctx.Done() 
    log.Println("Context Cancelled: Unsubcribing from ",topic)
    m.UnSubscribe(ctx,topic)
  }()
  m.subcribers[topic]=append(m.subcribers[topic],handler )
  return nil
}

func (m *memoryBackend)UnSubscribe(ctx context.Context,topic string)error{ 
  m.mu.RLock() 
  defer m.mu.RUnlock()
  delete(m.subcribers,topic)
  return nil
}

func (m *memoryBackend)SubscribeUser(ctx context.Context,userId,topic string, handler func(pubsub.Message) )error{

return nil
}

//no close method for inMemory subcribed messages
func (m *memoryBackend)Close()error{
   return nil
}
