package main

/// joinic
import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"time"
)

func enforceHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		clientID := c.Request().Header.Get("Client-ID")
		if clientID != "wju" {
			return c.String(403, "Client not found")
		}
		clientAccessKey := c.Request().Header.Get("Client-Key")
		if clientAccessKey != "secret" {
			return c.String(403, "Invalid Key1")
		}
		contentType := c.Request().Header.Get("Content-Type")
		if contentType != "application/json" {
			return c.String(403, "Invalid Key2")
		}
		return next(c)
	}
}

func myLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		var clog []byte
		t1 := time.Now()
		if err = next(c); err != nil {
			return err
		}
		diff := time.Now().Sub(t1)
		req := c.Request()
		res := c.Response()
		data := c.Get("filter")
		if data != nil {
			t := reflect.ValueOf(data)
			e := t.Elem()
			m := make(map[string]interface{})
			for i := 0; i < e.NumField(); i++ {
				mValue := e.Field(i)
				mType := e.Type().Field(i)
				fmt.Println(mValue, "Value", mType, "Type", e, "Elem")
				if _, ok := mType.Tag.Lookup("logging"); ok {
					m[mType.Name] = mValue.Interface()
				}
			}
			clog, _ = json.Marshal(m)
		}

		fmt.Println(
			t1.Format(time.RFC3339),
			diff,
			c.RealIP(),
			req.Method,
			req.Host,
			req.Header.Get(echo.HeaderContentLength),
			res.Size,
			string(clog),
		)
		return nil
	}
}
func main() {
	e := echo.New()

	e.Use(enforceHandler)
	e.Use(myLogMiddleware)

	e.GET("/", func(c echo.Context) error {
		time.Sleep(20 * time.Second)
		i := Info{Name: "PPEE", Email: "easndq@eoplqwr.gom", UserID: "IDIDIDIDIDID"}
		i.Filter(c)
		return c.JSON(http.StatusOK, &Response{200, "Pong", i})
	})

	e.Logger.Fatal(e.Start(":3000"))

}

type Logfilter struct {
	Level   string `json:"level"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	UserID  string `json:"userID"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	User interface{} `json:"user"`
}

type Info struct {
	UserID string `json:"user_id" logging:"true"`
	Email  string `json:"email" logging:"true"`
	Name   string `json:"name"`
}

func (i *Info) Filter(c echo.Context) {
	c.Set("filter", i)
}
