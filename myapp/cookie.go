package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func writeCookie(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = "jon"
	cookie.Expires = time.Now().Add(24 * time.Second)
	c.SetCookie(cookie)
	return c.String(http.StatusOK, "write a cookie")
}

func readCookie(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	fmt.Println(cookie.Name, cookie.Value)
	return c.String(http.StatusOK, "read a cookie")
}

func readAllCookies(c echo.Context) error {
	for _, cookie := range c.Cookies() {
		fmt.Println(cookie.Name, cookie.Value)
	}
	return c.String(http.StatusOK, " read all cookies")
}

func main() {
	e := echo.New()

	e.GET("/wcookie", writeCookie)
	e.GET("/rcookie", readCookie)
	e.GET("/racookie", readAllCookies)

	e.Logger.Fatal(e.Start(":1323"))
}
