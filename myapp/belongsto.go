package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// person belongs to company
type Person struct {
	gorm.Model
	Name      string
	CompanyID string  `gorm:"foreignKey"`
	Company   Company `gorm:"references:Name"`
}
type Company struct {
	Name   string `gorm:"primaryKey;size:100;index"` //ID가 아닌 컬럼으로 foreignKey를 설정하게 되면, 크기와 index를 명시해줘야한다.
	Region string
}

func main() {
	dsn := "root:./wjson./@tcp(localhost:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local" // username:password/@tcp(host:port)/database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Db 연결에 실패하였습니다.")
	}

	db.AutoMigrate(&Company{}, &Person{})

	//db.Debug().Create(&Person{Name: "Kim", Company: Company{Name: "SK", Region: "Seoul"}})
	//db.Debug().Create(&Person{Name: "Lee", Company: Company{Name: "KT", Region: "Busan"}})
	//db.Debug().Create(&Person{Name: "Park", Company: Company{Name: "LG", Region: "Suwon"}})
	//db.Debug().Create(&Person{Name: "Kang", Company: Company{Name: "Plea", Region: "Seoul"}})

	var people []Person
	//db.Debug().Where("name Like ?", "%ee%").Preload("Company", "name Like ?", "%K%").Find(&people)

	//db.Debug().Preload("Company", "name = ?", "KT").Where("name Like ?", "k%").Find(&people) // preload("하위 구조체 명 *대소문자 구분").where("find 뒤에 붙는 구조체에게 적용할 조건")
	// Preload에 조건을 걸었을 때, Company에 해당하는 조건만 가져와서 people과 매핑해준다.
	// where에 있는 조건이랑 동시 적용이 아니여서, 두가지 조건을 모두 일치하는것이 출력되는게 아니라 people에는 where 절에 있는 조건대로 데이터가 들어가고 people의 company에는 preload의 조건에 맞으면 들어가고 아니면 없다

	//db.Debug().Joins("Company").Where("people.name Like ? and Company.name = ?", "K%", "Plea").Find(&people) // => Joins의 파라미터로는 테이블이름 말고 연결된 구조체 이름

	for _, person := range people {
		fmt.Println(person.Company.Region, person.CompanyID, person.Name, person.ID, person.Company)
	}
}
