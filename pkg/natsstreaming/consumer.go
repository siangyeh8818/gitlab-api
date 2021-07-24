package natsstreaming

import (
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
)

func Consumermessage(messages <-chan *message.Message) {
	for msg := range messages {
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
