package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type CreditCard struct {
	gorm.Model
	Number string
	//User2Id uint `gorm:"primaryKey"`
}

type User2 struct {
	gorm.Model
	Name         string
	CreditCardID uint
	CreditCard   CreditCard `gorm:"foreignKey:credit_card_id;reference:id"`
}

type result struct {
	Id     string
	Number string
	Name   string
}

func (u User2) tableName() string {
	return "user2"
}

func (c *CreditCard) tableName() string {
	return "credit_cards"
}

func (u *User2) BeforeDelete(tx *gorm.DB) (err error) {
	if u.Name == "Silver" {
		fmt.Println("silver user not allowed to delete")
		return errors.New("DELETE REJECT")
	}
	return
}

func main() {
	dsn := "root:./wjson./@tcp(localhost:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local" // username:password/@tcp(host:port)/database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Db 연결에 실패하였습니다.")
	}

	//db.AutoMigrate(&User2{}, &CreditCard{})
	//
	//db.Create(&User2{Name: "Silver", CreditCard: CreditCard{Number: "123-4231"}})
	//
	//db.Create(&User2{Name: "Gold"})
	//db.Create(&User2{Name: "Sky", CreditCard: CreditCard{Number: "987-666"}})

	var userTest User2
	//db.Debug().Raw("select * from user2 where id = 13").Find(&userTest)
	//	db.Debug().Where("id = ?", 13).First(&userTest) //hooks 메서드에서 객체를 초기화시켜주지 않기 때문에, 메서드 전에 객체에 데이터를 담고 메서드 실행한다.
	//db.Debug().Delete(&userTest)

	//db.Exec("insert into credit_cards values(7,now(), now(), null,'100-200')")  // sql 문구를 그대로 실행
	//db.Exec("insert into user2 values (14, now(), now(), null, 'Black',7)")

	//	db.Debug().Unscoped().Where("id=13").Find(&userTest)
	//	fmt.Println("#삭제된 user", userTest)

	// # user 하위의 creditCard 구조체를 만들어서 user과 연결된 foreignkey를 수정
	//var c CreditCard
	//c.Number = "888-555"
	//db.Create(&c)
	//db.Model(&userTest).Where("name = ?", "Blue").Update("credit_card_id", &c.ID)

	//db.Where(User2{Name: "Navy"}).Attrs(User2{CreditCard: CreditCard{Number: "5546-7212"}}).FirstOrInit(&userTest) // where 조건으로 찾는 데이터가 없을 경우 구조체에 다른 값을 넣어 반환, db에 등록되지는 않는다

	//db.Where(User2{Name: "Navy"}).Attrs(User2{CreditCard: CreditCard{Number: "5546-7212"}}).FirstOrCreate(&userTest) // where 조건으로 찾는 데이터가 없을 경우 db에 등록하고 해당 데이터를 구조체에 넣어서 반환, 데이터가 있을 경우 속성은 무시됌
	//fmt.Println(userTest.Name, userTest.CreditCard.Number, userTest.ID, userTest.CreditCard)

	//db.Preload("CreditCard").Where("name =?", "Navy").Find(&userTest) //  Preload 사용 안하면 user2안의 creditcard에는 내용 안들어감
	//fmt.Println(userTest.Name, userTest.CreditCard.Number, userTest.ID, userTest.CreditCard)

	start := time.Now()
	// # 하기 2개의 메서드는 같은 기능
	db.Debug().Select("user2.id", "user2.name", "credit_cards.number AS CreditCard__number", "credit_cards.id AS CreditCard__id").Where("name = ?", "Silver").Joins("inner join credit_cards on user2.credit_card_id = credit_cards.id").Take(&userTest)
	//17.5909ms    join으로 가져올 때 join한 테이블의 필드는 alias 사용하는데 구조체명__변수명으로 해야 인식가능하다.

	//db.Debug().Preload("CreditCard").Where("name = ?", "Silver").Find(&userTest) //17.1263ms    // preload 사용 시 foreign key가 없는 데이터면 실행이 안되고 그 이후 쿼리문만 작동되고 다음으로 넘어감

	fmt.Println(userTest.ID, userTest.Name, userTest.CreditCard.Number, userTest.CreditCard.ID, userTest.CreditCardID)

	fin := time.Since(start)

	fmt.Println(fin)

}
