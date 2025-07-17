package main

import (
	"context"
	"log"
	"net/smtp"
	"time"

	"github.com/MohdAhzan/goPubSub/backend/memory"
	"github.com/MohdAhzan/goPubSub/pubsub"
)

func main(){

  ctx :=context.Background() 
    
  inMemoryPS :=memory.New(ctx)
  err:=inMemoryPS.Subscribe(ctx,"News",func(msg pubsub.Message){
    log.Println("General Handler",string(msg.Payload))

  })

  if err!=nil{
    log.Fatal(err)
  }


  err=inMemoryPS.Subscribe(ctx,"orders",func(msg pubsub.Message){
    log.Println("General Handler",string(msg.Payload),msg.UserID)

  })
  if err!=nil{
    log.Fatal(err)
  }

  err=inMemoryPS.Publish(ctx,pubsub.Message{
    ID: "Uuid123124",
    Topic: "News",
     Payload: []byte("random news"),
    Timestamp: time.Now(),
  })
  if err!=nil{
    log.Fatal(err)
  }

  err=inMemoryPS.Publish(ctx,pubsub.Message{
    ID: "Uuid123124",
    Topic: "orders",
    UserID: "1",
     Payload: []byte("order placed"),
    Timestamp: time.Now(),
  })
  if err!=nil{
    log.Fatal(err)
  }

    select{}
}

func sendmail(Msg pubsub.Message)error{

	TO := []string{"nejipkm@gmail.com"}

  SMTP_USERNAME:="gsirkahzanpkm@gmail.com"
  SMTP_PASSWORD:="wxgm xxxm otzg pbes"
  SMTP_HOST:="smtp.gmail.com"
  SMTP_PORT:="587"
//
	//setup authentication
	auth := smtp.PlainAuth("", SMTP_USERNAME, SMTP_PASSWORD, SMTP_HOST)

	//message body
	msg := []byte("To: " + TO[0] + "\r\n" +
		"Subject: " + Msg.Topic+ "\r\n" +
		"\r\n" +
		string(Msg.Payload) + "\r\n"+
		Msg.Timestamp.String()+ "\r\n")
	//send mail to recipient
	err := smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, auth, SMTP_USERNAME, TO, msg)
	if err != nil {
		return err
	}
	return nil


}
