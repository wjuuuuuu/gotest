package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	e := echo.New()
	e.Use(middleware.RequestID())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, c.Response().Header().Get(echo.HeaderXRequestID)) // 요청에 대한 고유 ID
	})

	e.Logger.Fatal(e.Start(":1324"))

}
