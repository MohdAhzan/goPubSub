package pubsub

import "time"

type Message struct{
  ID string
  Topic string
  UserID string // use for user-level delivery
  Payload []byte
  Timestamp time.Time

}


