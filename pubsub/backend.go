package pubsub

import "context"

type Backend interface{
  Publish(ctx context.Context,msg Message)error
  Subscribe(ctx context.Context,topic string,handler func(Message))error
  SubscribeUser(ctx context.Context,userId,topic string,handler func(Message))error
  UnSubscribe(ctx context.Context,topic string)error
  Close()error
}
