package memory

import (
	"context"
	"log"
	"sync"

	"github.com/MohdAhzan/goPubSub/pubsub"
)


type memoryBackend struct{

  mu                sync.RWMutex
  subcribers        map[string][]func(pubsub.Message)
  userHandlers      map[string]map[string][]func(pubsub.Message)
  queues            map[string]chan pubsub.Message
  workers           map[string]context.CancelFunc
  ctx               context.Context
  cancel            context.CancelFunc

}

func New(parent context.Context)pubsub.Backend{

  ctx,cancel:=context.WithCancel(parent)

  return &memoryBackend{
    subcribers: make(map[string][]func(pubsub.Message)),
    userHandlers: make(map[string]map[string][]func(pubsub.Message)),
    queues: make(map[string]chan pubsub.Message),
    workers: make(map[string]context.CancelFunc),
    ctx: ctx,
    cancel: cancel,
  }
}

func (m *memoryBackend)Publish(ctx context.Context,msg pubsub.Message)error{

  m.mu.RLock()
  queue,ok:=m.queues[msg.Topic]
  log.Println("QUEUE CHECKING ")
  m.mu.RUnlock()

  if !ok{
  log.Println("NO QUEUE STARTING WORKER")
    m.startWorker(msg.Topic)
    m.mu.RLock()
    queue=m.queues[msg.Topic]
    m.mu.RUnlock()
  }

  select {
  case queue <- msg:
    return nil
  case <-ctx.Done():
    return ctx.Err()
  }
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

  m.mu.RLock()
  defer m.mu.RUnlock()

  if m.userHandlers[userId]==nil{

    m.userHandlers[userId]=make(map[string][]func(pubsub.Message))
  }
  m.userHandlers[userId][topic]=append(m.userHandlers[userId][topic], handler)

  return nil
}

//no close method for inMemory subcribed messages
func (m *memoryBackend)Close()error{
  return nil
}
