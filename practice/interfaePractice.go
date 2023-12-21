package main

import "fmt"

type Animal interface {
	Bark()
	Eat()
	Run()
}

type Cats struct {
	Name string
	*Life
}

func (c *Cats) Bark() {
	fmt.Println(c.Name + " BARK")
}
func (c *Cats) Eat() {
	fmt.Println(c.Name + " EAT")
}

type Life struct {
	Number int
}

func (l *Life) Run() {
	fmt.Println(l.Number, "Run")
}

type Elephant struct {
	Name string
}

func (e *Elephant) Bark() {
	fmt.Println(e.Name + " BARK")
}
func (e *Elephant) Eat() {
	fmt.Println(e.Name + " EAT")
}
func main() {
	cat := &Cats{"FOO", &Life{Number: 12}}
	var a Animal
	a = cat
	a.Bark()
	a.Eat()
	a.Run()

	//var b Animal
	// b = &Elephant{Name: "BOO"} // RUN()을 가지고 있지 않아서 Elephant 구조체는 Animal 인터페이스가 아니다
	//b.Bark()
	//b.Eat()
	//b.Run()

}
