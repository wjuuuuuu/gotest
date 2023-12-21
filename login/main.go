package main

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	db2 "login/db"
	"login/helper"
	"login/nats"
	"login/user"
	"net/http"
	"os"
	"time"
)

type Context struct {
	echo.Context
	DB *gorm.DB
}

type loginIDPWRequest struct {
	ID string `json:"id"`
	PW string `json:"pw"`
}

func main() {
	db, err := db2.DBConnect()
	if err != nil {
		os.Exit(1)
	}

	nats.NatsConnect()
	nats.NatsSubscribe()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut},
		AllowHeaders: []string{"*"},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	grp := e.Group("/user")

	grp.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := &Context{c, db}
			return next(ctx)
		}
	})

	grp.GET("/", CreateTableUser)
	grp.POST("/join", InsertTableUser)
	grp.POST("/login", Login)
	grp.POST("/password", ModifyPasswordUser)
	grp.PUT("/status/:id", ModifyStatusUser)

	e.Logger.Fatal(e.Start(":3000"))

}

func CreateTableUser(c echo.Context) error {
	db := c.(*Context).DB
	err := db.AutoMigrate(&user.User{})
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "테이블 만들기 실패",
			"data":    nil,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "테이블 만들기 성공",
		"data":    nil,
	})
}

func InsertTableUser(c echo.Context) error {
	db := c.(*Context).DB
	newUser := &user.User{}

	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "bad request",
			"data":    nil,
		})
	}
	if newUser.FindDuplicateDateUser(newUser.Ctn, db) != 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "이미 있는 계정입니다.",
			"data":    nil,
		})
	}
	hashPw, err := helper.HashPassword(newUser.Password)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err,
			"data":    nil,
		})
	}
	newUser.Password = hashPw

	if err := db.Create(&newUser).Error; err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "계정을 생성할 수 없습니다.",
			"data":    nil,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "계정 생성 완료",
		"data":    newUser,
	})
}

func Login(c echo.Context) error {
	db := c.(*Context).DB
	idPwTmp := &loginIDPWRequest{}

	if err := c.Bind(&idPwTmp); err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "bad request",
			"data":    nil,
		})
	}

	if idPwTmp.ID == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "아이디를 입력하세요",
			"data":    nil,
		})
	}
	if idPwTmp.PW == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "비밀번호를 입력하세요",
			"data":    nil,
		})
	}

	loginUser := &user.User{}
	if err := db.Where("id = ?", idPwTmp.ID).Take(&loginUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "없는 ID 입니다.",
				"data":    nil,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "다시 시도해주세요",
			"data":    nil,
		})
	}

	pwCheck := helper.CheckPasswordHash(loginUser.Password, idPwTmp.PW)
	if pwCheck == false {
		failCount, err := loginUser.AddFailCount(db)
		if err != nil {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": "다시 시도해주세요",
				"data":    nil,
			})
		}
		if failCount >= 5 {
			user.ModifyStatus(idPwTmp.ID, "2", db)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "로그인 불가 계정, 관리자에게 문의",
				"data":    nil,
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "아이디 비밀번호 불일치",
			"data":    nil,
		})
	}

	if loginUser.Status != user.Active {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "로그인 불가 계정, 관리자에게 문의",
			"data":    nil,
		})
	}

	accessToken, err := helper.CreateJWT(idPwTmp.ID)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "다시 시도해주세요",
			"data":    nil,
		})
	}

	cookie := new(http.Cookie)
	cookie.Name = "access-token"
	cookie.Value = accessToken
	cookie.Expires = time.Now().Add(time.Hour * 1)
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "로그인 성공",
		"data":    loginUser,
	})

}

func ModifyStatusUser(c echo.Context) error {
	db := c.(*Context).DB

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "bad request",
			"data":    nil,
		})
	}

	newStatus := c.QueryParam("status")
	if newStatus == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "bad request.",
			"data":    nil,
		})
	}

	err := user.ModifyStatus(id, newStatus, db)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "Fail to update status",
			"data":    nil,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "상태 변경 완료",
		"data":    nil,
	})
}

func ModifyPasswordUser(c echo.Context) error {
	db := c.(*Context).DB
	idPwTmp := &loginIDPWRequest{}

	if err := c.Bind(&idPwTmp); err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "bad request",
			"data":    nil,
		})
	}

	if idPwTmp.ID == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "아이디를 입력하세요",
			"data":    nil,
		})
	}
	if idPwTmp.PW == "" {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "비밀번호를 입력하세요",
			"data":    nil,
		})
	}

	hashPw, err := helper.HashPassword(idPwTmp.PW)
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": err,
			"data":    nil,
		})
	}
	tmpUser := user.User{}
	err = tmpUser.ModifyPassword(idPwTmp.ID, hashPw, db)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "없는 아이디",
				"data":    nil,
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "서버에러",
			"data":    nil,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    http.StatusOK,
		"message": "비밀번호 변경",
		"data":    tmpUser,
	})

}
