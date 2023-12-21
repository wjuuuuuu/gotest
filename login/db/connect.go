package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DBConnect() (*gorm.DB, error) {
	dsn := "root:./wjson./@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local" // username:password/@tcp(host:port)/database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Db 연결에 실패하였습니다.")
		return nil, err
	}
	return db, nil
}
