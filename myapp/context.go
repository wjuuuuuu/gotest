package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

type CustomContext struct {
	echo.Context
}

func (c *CustomContext) Foo() {
	println("foo")
}

func (c *CustomContext) Bar() {
	println("bar")
}

func main() {
	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c}
			fmt.Println("cc", cc.Context)
			fmt.Println("next", next)
			return next(cc)
		}
	})
	e.GET("/", func(c echo.Context) error {
		fmt.Println("c", c)
		cc := c.(*CustomContext)
		cc.Foo()
		cc.Bar()
		return cc.String(200, "OK")
	})
	e.Logger.Fatal(e.Start(":1324"))
}
