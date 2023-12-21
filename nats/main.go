package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func main() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		os.Exit(1)
	}
	go subs(nc)

	sub, _ := nc.SubscribeSync("greet.*")
	msg, _ := sub.NextMsg(1 * time.Second)

	nc.Publish("greet.joe", []byte("eeee"))
	fmt.Println("subscribed after a publish...")
	fmt.Printf("msg is nil? %v\n", msg == nil)

	nc.Publish("greet.qqqq", []byte("1111"))
	nc.Publish("greet.ddd", []byte("he22llo"))
	nc.Publish("greet.wwwq", []byte("he33llo"))

	msg, _ = sub.NextMsg(2 * time.Second)
	fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)

	msg, _ = sub.NextMsg(time.Second)
	fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)

	msg, _ = sub.NextMsg(time.Second)
	fmt.Printf("msg data: %q on subject %q\n", string(msg.Data), msg.Subject)

	wg.Add(1)
	wg.Wait()

}

func subs(conn *nats.Conn) {
	conn.Subscribe("greet.*", func(m *nats.Msg) {
		fmt.Println("비동기", string(m.Data), string(m.Subject))
	})
	time.Sleep(10 * time.Second)
	wg.Done()
}
