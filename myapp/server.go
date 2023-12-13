package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Car struct {
	name  string `query:"name"`
	price int    `query:"price"`
}

func getUsers(c echo.Context) error {
	users := []User{{1, "ryan"}, {2, "john"}}
	return c.JSON(http.StatusOK, users)
}
func postUsers(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		c.String(http.StatusBadRequest, "bad request")
	}
	fmt.Println(user)
	return c.JSON(http.StatusCreated, user)
}

func main() {
	e := echo.New()
	e.GET("/users", getUsers)
	e.POST("/users", postUsers)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.GET(("/car"), func(c echo.Context) error {
		fmt.Println(c.QueryParam("name"), "name")
		fmt.Println(c.QueryParam("price"), "price")
		var car Car
		car.name = c.QueryParam("name")
		price := c.QueryParam("price")
		intprice, _ := strconv.Atoi(price)
		car.price = intprice

		return c.String(http.StatusOK, car.name+price)

	})
	e.Logger.Fatal(e.Start(":1324"))
}
