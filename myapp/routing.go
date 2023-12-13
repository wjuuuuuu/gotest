package main

import (
	"crypto/subtle"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func hello(c echo.Context) error {
	param := c.Param("param")
	return c.String(http.StatusOK, "HEllo,"+param)
}

func main() {
	e := echo.New()

	g := e.Group("/admin")
	//g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	//	if username == "joe" && password == "secret" {
	//		return true, nil
	//	}
	//	return false, nil
	//}))

	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte("chris")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("12345")) == 1 {
			return true, nil
		}
		return false, nil
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(echoprometheus.NewMiddleware("myapp"))
	e.GET("/metrics", echoprometheus.NewHandler())
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello")
	})

	e.GET("/admin/:param", hello)

	e.Logger.Fatal(e.Start(":1324"))

}
