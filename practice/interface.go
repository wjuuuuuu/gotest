package main

import "fmt"

type Speaker interface {
	Speak()
}

type Dog struct {
	Name       string
	isMammal   bool
	PackFactor int
}
type Cat struct {
	Name        string
	isMammal    bool
	ClimbFactor int
}

func (d Dog) Speak() {
	fmt.Println("Woof! My name is", d.Name, "it is ", d.isMammal, "I am mammal with pack factor of", d.PackFactor)
}

func (c Cat) Speak() {
	fmt.Println("Meow! My name is", c.Name, "it is ", c.isMammal, "I am mammal with pack factor of", c.ClimbFactor)
}
func main() {
	speakers := []Speaker{
		Dog{Name: "Boo", isMammal: true, PackFactor: 3},
		Cat{Name: "Foo", isMammal: true, ClimbFactor: 2},
	}

	for _, s := range speakers {
		s.Speak()
	}
}
