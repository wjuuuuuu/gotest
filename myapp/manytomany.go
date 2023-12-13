package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

const (
	VIP = iota
	GOLD
	Silver
	Green
)

type Grade int

var gradeNames = map[Grade]string{
	VIP:    "VIP",
	GOLD:   "GOLD",
	Silver: "Silver",
	Green:  "Green",
}

func stringToGrade(s string) (Grade, error) {
	for grade, name := range gradeNames {
		if name == s {
			return grade, nil
		}
	}
	return 0, fmt.Errorf("Invalid grade: %s", s)
}

type Person2 struct {
	gorm.Model
	Name  string
	Grade Grade
}
type Product2 struct {
	gorm.Model
	Name  string
	Price int
}
type Order struct {
	PersonID  int `gorm:"primaryKey;foreignKey"`
	ProductID int `gorm:"primaryKey;foreignKey"`
	Count     int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
	Product   Product2 `gorm:"references:id;"`
	Person    Person2  `gorm:"references:id;"`
}

func (p Person2) TableName() string {
	return "person"
}
func (p Product2) TableName() string {
	return "product"
}
func (o Order) TableName() string {
	return "order"
}

type Amount struct {
	Name string
	Tot  int
}

func findPersonByName(c echo.Context, db *gorm.DB) error {
	userName := c.QueryParam("userName")
	var person []Person2
	if err := db.Model(&Person2{}).Where("name Like ?", "%"+userName+"%").Find(&person).Error; err != nil {
		return err
	}
	if len(person) == 0 {
		c.String(http.StatusOK, "not matched person")
		return nil
	}
	for _, v := range person {
		strGrade := gradeNames[v.Grade]
		c.String(http.StatusOK, v.Name+" "+strGrade+"\n")
	}
	return nil
}

func createPerson(c echo.Context, db *gorm.DB) error {
	userName := c.QueryParam("userName")
	grade := c.QueryParam("grade")
	intGrade, err := stringToGrade(grade)
	if err != nil {
		return err
	}
	var person = Person2{Name: userName, Grade: intGrade}
	err = db.Debug().Create(&person).Error
	if err != nil {
		return err
	}
	c.String(http.StatusCreated, "CREATE PERSON")
	return nil
}

func updatePerson(c echo.Context, db *gorm.DB) error {
	id := c.QueryParam("id")
	changeName := c.QueryParam("name")
	var p Person2
	db.Where("id = ?", id).First(&p)

	if p.ID == 0 {
		err := c.String(http.StatusNotFound, "해당 id의 회원은 존재하지 않습니다.")
		if err != nil {
			return err
		}
	}
	if changeName == "" {
		err := c.String(http.StatusNotFound, "변경할 이름을 입력해주세요")
		return err
	}
	if err := db.Model(&p).Where("id = ?", id).Update("name", changeName).Error; err != nil {
		fmt.Println("update 실패")
		return err
	}
	err := c.String(http.StatusOK, "UPDATE Person")
	if err != nil {
		return err
	}
	return nil
}

func deletePerson(c echo.Context, db *gorm.DB) error {
	tx := db.Begin()
	id := c.QueryParam("id")
	var p Person2
	db.Where("id = ?", id).First(&p)

	if p.ID == 0 {
		err := c.String(http.StatusNotFound, "해당 id의 회원은 존재하지 않습니다.")
		return err
	}

	tx.Where("id = ?", id).Delete(&p)

	if p.Name == "Kim" {
		tx.Rollback()
		err := c.String(http.StatusBadRequest, "Kim은 삭제 불가")
		return err
	}
	tx.Commit()
	err := c.String(http.StatusOK, "DELETE COMPLETE")
	if err != nil {
		return err
	}
	return nil
}

func orderProduct(c echo.Context, db *gorm.DB) error {
	userName := c.QueryParam("userName")
	proName := c.QueryParam("proName")
	count := c.QueryParam("count")

	intcnt, err := strconv.Atoi(count)
	if err != nil || intcnt <= 0 {
		c.String(http.StatusBadRequest, "수량이 잘못됐습니다.")
		return errors.New("수량 잘못")
	}

	var p Person2
	var pro Product2
	db.Where("name = ?", userName).First(&p)
	db.Where("name = ?", proName).First(&pro)

	if p.ID == 0 {
		c.String(http.StatusBadRequest, "회원이 아닙니다")
		return errors.New("id 없음")
	}
	if pro.ID == 0 {
		c.String(http.StatusBadRequest, "상품이 아닙니다.")
		return errors.New("id 없음")
	}

	err = db.Create(&Order{Person: p, Product: pro, Count: intcnt, CreatedAt: time.Now()}).Error
	if err != nil {
		c.String(http.StatusUnauthorized, "구매실패")
		return err
	}
	tot := strconv.Itoa(pro.Price * intcnt)
	c.String(http.StatusOK, p.Name+"님 "+pro.Name+" 구매 : "+tot+"원")

	return nil
}
func main() {

	e := echo.New()

	dsn := "root:./wjson./@tcp(localhost:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local" // username:password/@tcp(host:port)/database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Db 연결에 실패하였습니다.")
	}

	db.Migrator().CreateIndex(&Person2{}, "Name")
	bools := db.Migrator().HasIndex(&Person2{}, "Name")
	fmt.Println(bools)
	db.AutoMigrate(&Person2{}, &Product2{}, &Order{})

	//tx := db.Begin()
	//person := []Person2{
	//	{Name: "Kim", Grade: VIP},
	//	{Name: "Lim", Grade: GOLD},
	//	{Name: "Zia", Grade: Green},
	//	{Name: "John", Grade: Silver},
	//}
	//
	//products := []Product2{
	//	{Name: "apple", Price: 600},
	//	{Name: "banana", Price: 400},
	//	{Name: "donut", Price: 200},
	//	{Name: "carrot", Price: 100},
	//	{Name: "ice-cream", Price: 1400},
	//	{Name: "milk", Price: 900},
	//}
	//db.Create(&person)
	//db.Create(&products)
	//
	//orders := []Order{
	//	{PersonID: 1, ProductID: 1, CreatedAt: time.Now(), Count: 10},
	//	{PersonID: 1, ProductID: 5, CreatedAt: time.Now(), Count: 3},
	//	{PersonID: 1, ProductID: 6, CreatedAt: time.Now(), Count: 2},
	//	{PersonID: 2, ProductID: 3, CreatedAt: time.Now(), Count: 7},
	//	{PersonID: 2, ProductID: 2, CreatedAt: time.Now(), Count: 5},
	//	{PersonID: 2, ProductID: 1, CreatedAt: time.Now(), Count: 12},
	//	{PersonID: 3, ProductID: 4, CreatedAt: time.Now(), Count: 8},
	//	{PersonID: 3, ProductID: 5, CreatedAt: time.Now(), Count: 10},
	//	{PersonID: 4, ProductID: 2, CreatedAt: time.Now(), Count: 6},
	//	{PersonID: 4, ProductID: 3, CreatedAt: time.Now(), Count: 7},
	//	{PersonID: 4, ProductID: 6, CreatedAt: time.Now(), Count: 5},
	//}
	//if err := tx.Create(&orders).Error; err != nil {
	//	fmt.Println(err)
	//	tx.Rollback()
	//	return
	//} else {
	//	tx.Commit()
	//}

	// # KIM이 구매한 내역 가져오기, Preload 사용하면 연관된 전체 데이터를 전부 가져올 때는 유리한듯.
	//var orders []Order
	//db.Debug().Preload("Product").Preload("Person").Find(&orders) // preload 에 조건을 걸어도 연관된 데이터를 다 가져오기 때문에 find로 찾으면 조건에 해당되지 않는 내용도 불러온다..
	//for _, val := range orders {
	//	if val.Person.Name == "Kim" {
	//		fmt.Println(val.Person.Name, val.Product.Name, val.Product.Price*val.Count, val.CreatedAt)
	//	}
	//}
	//print("\n\n")

	// # person name별 총 구매 금액 구하기
	//var a []Amount
	//db.Debug().Raw("select e.user as name, sum(e.purchase) as tot from (  select pe.name as user, pro.name, pro.price *o.count as purchase ,o.created_at from `order` o inner join person pe  on o.person_id = pe.id inner join product pro  on o.product_id = pro.id)e group by e.user").Find(&a)
	// # 섭쿼리 사용
	//subquery := db.Table("`order` o").
	//	Select("pe.name as name, (pro.price * o.count) as tot").
	//	Joins("INNER JOIN Person pe ON o.person_id = pe.id").
	//	Joins("INNER JOIN Product pro ON o.product_id = pro.id")
	//
	//db.Debug().Table("(?) as e", subquery).Select("name, sum(tot) as tot").Group("name").Find(&a) // subquery로 넣을 때 table alias 해줘야함
	//for _, am := range a {
	//	fmt.Println(am.Name, am.Tot)
	//}

	//var p Product2
	//queryDB := db.Debug().Where("price >  ?", 500).Session(&gorm.Session{}) // session없으면 계속 where 절 추가됌
	//queryDB.Where("price < ?", 1000).First(&p)
	//fmt.Println(p)
	//queryDB.Where("name Like ?", "%a%").First(&p)
	//fmt.Println(p)

	// 트랜잭션 save point
	//tx.Create(&Person2{Name: "test2", Grade: Green})
	//tx.SavePoint("S1") // save point 1
	//
	//tx.Create(&Person2{Name: "test3", Grade: VIP})
	//tx.SavePoint("S2") // save point2
	//
	//tx.RollbackTo("S1") // savepoint1로 롤백
	//tx.Commit()         // commit -> savepoint1 까지만 commit됌

	e.GET("/find", func(c echo.Context) error {
		return findPersonByName(c, db)
	})

	e.POST("/create", func(c echo.Context) error {
		return createPerson(c, db)
	})

	e.POST("/update", func(c echo.Context) error {
		return updatePerson(c, db)
	})

	e.POST("/delete", func(c echo.Context) error {
		return deletePerson(c, db)
	})

	e.POST("/order", func(c echo.Context) error {
		return orderProduct(c, db)
	})

	e.Logger.Fatal(e.Start(":1324"))
}
