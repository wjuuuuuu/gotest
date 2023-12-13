package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	dsn := "root:./wjson./@tcp(localhost:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local" // username:password/@tcp(host:port)/database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Db 연결에 실패하였습니다.")
	}

	// 테이블 자동 생성
	db.AutoMigrate(&Product{})

	//생성
	//db.Create(&Product{Code: "D42", Price: 100})

	// 읽기
	var product Product
	db.Find(&product, "id = ?", "1")
	fmt.Println(product.ID, product.Code, product.Price)
	var product2 Product
	db.Find(&product2, 2)
	fmt.Println(product2.ID, product2.Code, product2.Price)
	//db.First(&product, 1)                 // primary key기준으로 product 찾기
	//db.First(&product, "code = ?", "D42") // code가 D42인 product 찾기

	// 수정 - product의 price를 200으로
	db.Model(&product2).Update("Price", 1000)
	// 수정 - 여러개의 필드를 수정하기
	db.Model(&product).Updates(Product{Price: 1200, Code: "F42"})

	var product3 Product
	db.First(&product3, "id = ?", "3")
	db.Model(&product3).Updates(Product{Price: 5000, Code: "E652"})

	var productT Product
	db.First(&productT, "price = ?", 100)
	var productS []Product
	db.Find(&productS, "price= ?", 100)
	fmt.Println(productT.ID)
	for _, v := range productS {
		fmt.Println(v.ID, "#")
	}
	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// 삭제 - product 삭제하기
	db.Delete(&product, 4)
}
