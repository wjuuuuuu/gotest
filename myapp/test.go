package main

import "fmt"

type notifier interface {
	notify()
}
type duration int

func (d *duration) notify() {
	fmt.Println("Sending Notification in", *d)
}

func main() {
	d := duration(42) // 인라인으로는 안됌 42가변수에 저장된 값이 아니기 때문에 주소값을 얻을 수 없음.
	d.notify()
}
