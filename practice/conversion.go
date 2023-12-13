package main

import "fmt"

type Mover interface {
	Move()
}
type Locker interface {
	Lock()
	Unlock()
}

type MoveLocker interface {
	Mover
	Locker
}

type bike struct{}

func (bike) Move() {
	fmt.Println("Moving the bike")
}
func (bike) Lock() {
	fmt.Println("Locking the bike")
}
func (bike) Unlock() {
	fmt.Println("Unlocking the bike")
}

func main() {
	var ml MoveLocker
	var m Mover

	ml = bike{}
	m = ml
	//ml = m 불가능 MoveLocker가 더 구체화된 개념

	b, ok := m.(bike)
	fmt.Println("Does m has value of bike?", ok)
	ml = b

	fmt.Println(ml)

}
