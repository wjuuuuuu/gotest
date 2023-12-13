package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func main() {

	e := echo.New()
	e.GET("/", trigger)

	e.Logger.Fatal(e.Start(":1324"))
}

func trigger(c echo.Context) error {
	ca := make(chan string, 1)
	r := c.Request()
	method := r.Method

	go func() {
		fmt.Println("METHOD: ", method)

		ca <- "HEY"
	}()

	select {

	case result := <-ca:
		return c.String(http.StatusOK, "Result: "+result)
	case <-time.After(20 * time.Second):
		return nil
	}

}
