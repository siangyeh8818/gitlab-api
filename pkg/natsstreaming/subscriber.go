package natsstreaming

import (
	"time"

	stan "github.com/nats-io/stan.go"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/siangyeh8818/gitlab.api/pkg/setting"
)

func SubScriber(clientid string) *nats.StreamingSubscriber {

	//起一個訂閱者的client
	subscriber, err := nats.NewStreamingSubscriber(
		nats.StreamingSubscriberConfig{
			ClusterID:        setting.NatsStreamingSetting.ClusterId,
			ClientID:         clientid,
			QueueGroup:       setting.NatsStreamingSetting.QueueGroup,
			DurableName:      "my-durable",
			SubscribersCount: 4, // how many goroutines should consume messages
			CloseTimeout:     time.Minute,
			AckWaitTimeout:   time.Second * 30,
			StanOptions: []stan.Option{
				stan.NatsURL(setting.NatsStreamingSetting.Address),
			},
			Unmarshaler: nats.GobMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}
	return subscriber
}
