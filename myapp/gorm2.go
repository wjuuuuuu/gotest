package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Student1 struct {
	gorm.Model
	Name  string
	Grade int
	Etc   string
}

func DbConnect() *gorm.DB {
	dsn := "root:./wjson./@tcp(localhost:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local" // username:password/@tcp(host:port)/database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Db 연결에 실패하였습니다.")
	}
	return db
}

func main() {

	db := DbConnect()
	db.AutoMigrate(&Student1{})

	students := []Student1{
		Student1{Name: "Paul", Grade: 3, Etc: "he is chief"},
		Student1{Name: "Julia", Grade: 2, Etc: "she has a book"},
		Student1{Name: "Tom", Grade: 1, Etc: "he looks like a teacher"},
	}

	studentA := Student1{Name: "Ping", Grade: 2}
	db.Select("Name", "grade").Create(&studentA)

	result := db.Create(&students)
	fmt.Println(result.Error)
	fmt.Println(result.RowsAffected)

	//db.Delete(&Student1{}, 5)

}
