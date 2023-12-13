package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Users struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}
type UsersDto struct {
	Name    string
	Email   string
	IsAdmin bool
}

func main() {
	e := echo.New()

	e.POST("/users", func(c echo.Context) error {
		u := new(Users)
		var age = int64(35)
		err := echo.QueryParamsBinder(c).String("name", &u.Name).String("email", &u.Email).Int64("age", &age).BindErrors()
		fmt.Println(err, "err")
		fmt.Println(age, "age")

		if err := c.Bind(u); err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		user := UsersDto{
			Name:    u.Name,
			Email:   u.Email,
			IsAdmin: false,
		}
		fmt.Println(user)
		return c.JSON(http.StatusOK, u)
	})

	// get 요청이 들어오면 채널이 생성되고 연결이 되면 고루틴으로 Hey를 받아서 출력하고 아니면  nil 반환
	e.GET("/:param", func(c echo.Context) error {
		param := c.Param("param")
		ca := make(chan string, 1) // To prevent this channel from blocking, size is set to 1.
		r := c.Request()
		method := r.Method

		go func() {
			// This function must not touch the Context.

			fmt.Printf("Method: %s\n", method)
			fmt.Println(r.URL)
			fmt.Println(c.Response().Header, " ## ", r.Header)

			// Do some long running operations...

			ca <- param
		}()

		select {
		case result := <-ca:
			return c.String(http.StatusOK, "Result: "+result)
		case <-c.Request().Context().Done(): // Check context.
			// If it reaches here, this means that context was canceled (a timeout was reached, etc.).
			return nil
		}
	})

	e.Logger.Fatal(e.Start(":1324"))
}
