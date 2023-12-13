package main

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Student struct {
	Name  string `query:"name" json:"name" validate:"required"`
	Grade int64  `query:"grade" json:"grade" validate:"required,min=2"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	//e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	//	return func(c echo.Context) error {
	//		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valiid credentials")
	//	}
	//})

	e.POST("/", func(c echo.Context) error {
		s1 := new(Student)
		err := c.Bind(s1)
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}
		fmt.Println(s1)
		return c.JSON(http.StatusOK, s1)
	})

	// validator 사용
	e.POST("/students", func(c echo.Context) (err error) {
		u := new(Student)
		if err = c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(u); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, u)
	})

	e.GET("/json", func(c echo.Context) error {
		c.Response().Before(func() {
			fmt.Println("before response")
		})
		s := &Student{
			Name:  "john",
			Grade: 2,
		}
		fmt.Println(s)

		return c.Redirect(http.StatusMovedPermanently, "/")
	})

	e.GET("/", func(c echo.Context) error {
		c.Response().After(func() {
			fmt.Println("after response")
		})
		fmt.Println("success")
		return c.String(http.StatusOK, "Redirect Success")
	})

	e.Logger.Fatal(e.Start(":1324"))

}
