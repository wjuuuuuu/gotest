package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"
)

//Xenia 라는 시스템은 데이테베이스를 가지고 있다.
//Pillar라는 또 다른 시스템은 프론트엔드를 가진 웹서버이며 Xenia를 이용한다.
//Pillar 역시 데이터베이스가 있다. Xenia의 데이터를 Pillar에 옮겨보자.

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Puller interface {
	Pull(d *Data) error
}
type Storer interface {
	Store(d *Data) error
}
type PullStorer interface {
	Puller
	Storer
}

type Data struct {
	Line string
}
type Xenia struct {
	Host    string
	Timeout time.Duration
}

func (*Xenia) Pull(d *Data) error {
	switch rand.Intn(10) {
	case 1, 9:
		return io.EOF
	case 5:
		return errors.New("Error reading data from Xenia")
	default:
		d.Line = "Data"
		fmt.Println("In:", d.Line)
		return nil
	}
}

type Pillar struct {
	Host    string
	Timeout time.Duration
}

func (*Pillar) Store(d *Data) error {
	fmt.Println("Out: ", d.Line)
	return nil
}

type System struct {
	Xenia
	Pillar
}

func pull(p Puller, data []Data) (int, error) {
	for i := range data {
		if err := p.Pull(&data[i]); err != nil {
			return i, err
		}
	}
	return len(data), nil
}

func store(s Storer, data []Data) (int, error) {
	for i := range data {
		if err := s.Store(&data[i]); err != nil {
			return i, err
		}
	}
	return len(data), nil
}

func Copy(ps PullStorer, batch int) error {
	data := make([]Data, batch)
	for {
		i, err := pull(ps, data)
		if i > 0 {
			if _, err := store(ps, data[:i]); err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}
}

func main() {
	sys := System{
		Xenia{Host: "localhost:8000", Timeout: time.Second},
		Pillar{Host: "localhost:9000", Timeout: time.Second},
	}
	if err := Copy(&sys, 3); err != io.EOF {
		fmt.Println(err)
	}
}
