package main

import "fmt"

type PubSub struct {
	host string
}

func New(host string) *PubSub {
	ps := PubSub{
		host: host,
	}
	return &ps
}
func (ps *PubSub) Publish(key string, v interface{}) error {
	fmt.Println("Actual PubSub: Publish")
	return nil
}
func (ps *PubSub) Subscribe(key string) error {
	fmt.Println("Actual PubSub: Subscribe")
	return nil
}

// 인터페이스를 만들어서..  Mock 객체를 활용
type Publisher interface {
	Publish(key string, v interface{}) error
	Subscribe(key string) error
}

type mock struct{}

func (m *mock) Publish(key string, v interface{}) error {
	fmt.Println("Mock PubSub: Publish")
	return nil
}
func (m *mock) Subscribe(key string) error {
	fmt.Println("Mock PubSub: Scribe")
	return nil
}

func main() {
	pubs := []Publisher{
		New("localhost"),
		&mock{},
	}

	for _, p := range pubs {
		p.Publish("key", "value")
		p.Subscribe("key")
	}
}
