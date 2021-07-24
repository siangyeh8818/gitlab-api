package natsstreaming

import (
	"fmt"
	"log"
	"os"
	"time"

	nats2 "github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"

	"github.com/siangyeh8818/gitlab.api/pkg/setting"
)

func ConnectNats() stan.Conn {

	opts := []nats2.Option{nats2.Timeout(10 * 60 * time.Second),
		nats2.MaxReconnects(50), nats2.ReconnectWait(10 * time.Second), nats2.ReconnectHandler(func(_ *nats2.Conn) {
			log.Println("nats client reconnected")
		})}

	nc, err := nats2.Connect(setting.NatsStreamingSetting.Address, opts...)

	if err != nil {
		log.Println("nats connect :", err)
	}
	defer nc.Close()

	sc, err := stan.Connect(setting.NatsStreamingSetting.ClusterId, os.Getenv("HOSTNAME"), stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Printf("Connection lost, reason: %v\n", reason)

		}))
	if err != nil {
		log.Println("Can't connect:", err)
		fmt.Printf("CMake sure a NATS Streaming Server is running at: %s", setting.NatsStreamingSetting.Address)

	}
	return sc

}
