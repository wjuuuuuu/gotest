package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Dog struct {
	ID   int `gorm:"primaryKey"`
	Name string
	Toy  []Toy `gorm:"foreignKey:DogID"` // toy는 dog에 속하게 dog의 pk를 foreignkey로
}
type Toy struct {
	Name  string
	Price int
	DogID int
}

func (d Dog) TableName() string {
	return "dog"
}
func (t Toy) TableName() string {
	return "toy"
}
func main() {
	dsn := "root:./wjson./@tcp(localhost:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local" // username:password/@tcp(host:port)/database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Db 연결에 실패하였습니다.")
	}

	db.AutoMigrate(Dog{}, Toy{})
	//dogs := []Dog{
	//	Dog{Name: "D1", Toy: []Toy{
	//		Toy{Name: "T1", Price: 1000},
	//		Toy{Name: "T2", Price: 1030},
	//		Toy{Name: "T3", Price: 1060},
	//	}},
	//	Dog{Name: "D2", Toy: []Toy{
	//		Toy{Name: "T4", Price: 2000},
	//		Toy{Name: "T5", Price: 2030},
	//		Toy{Name: "T6", Price: 2060},
	//	}},
	//}
	//db.Create(&dogs)

	var t []Toy
	var d Dog

	db.Where("id = ?", 5).First(&d) // 모델 데이터 선정

	c := db.Debug().Model(&d).Association("Toy").Count() // 모델 데이터, Toy 연관 갯수
	db.Debug().Model(&d).Association("Toy").Find(&t)     // 모델 데이터, Toy 연관있는 데이터
	fmt.Println(c)
	for _, v := range t {
		fmt.Println(v.DogID, v.Name, v.Price)
	}
	//db.Debug().Model(&d).Association("Toy").Append(&Toy{Name: "New-T", Price: 5000}) // 모델데이터와 연관되는 데이터 추가

	//newDog := Dog{
	//	Name: "D3", Toy: []Toy{
	//		{Name: "T1", Price: 7010},
	//		{Name: "T2", Price: 7020},
	//	},
	//}
	//db.Create(&newDog)
	//db.Debug().Omit(clause.Associations).Create(&newDog) //  dog와 연관된 테이블은 skip

	db.Debug().Model(&d).Association("Toy").Replace(&Toy{Name: "Test", Price: 1234, DogID: d.ID})

	//err = db.Debug().Model(&d).Association("Toy").Replace(&newToy)
	//if err != nil {
	//	return
	//}
	db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&d) // 변경된 사항을 반영해줌

	//var printDog []Dog
	//db.Preload("Toy", "price > ?", 2000).Where("name Like ?", "%2%").Find(&printDog)
	//for _, dog := range printDog {
	//	for _, toy := range dog.Toy {
	//		fmt.Println(dog.ID, dog.Name, toy.Name, toy.Price, toy.DogID)
	//	}
	//
	//}
}
