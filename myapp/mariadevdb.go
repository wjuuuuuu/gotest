package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type BizCustomer struct {
	gorm.Model
	UserInfo        string
	UserName        string
	CompanyName     string
	CustomerGroupId string
}

func main() {
	dns := "dxm:Cjswo.123@tcp(im.plea.kr:13306)/dxm?parseTime=True"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic("db connection fail")
	}

	var bizcustom []BizCustomer

	db.Table("biz_customer").Find(&bizcustom, "user_name = ?", "plea")

	for i, v := range bizcustom {
		fmt.Println(i, v.UserInfo, v.UserName, v.CompanyName, v.ID)
	}

	if len(bizcustom) > 1 {
		biz := bizcustom[0]
		db.Table("biz_customer").Model(&biz).Update("user_name", "플리test")
	}

	var biztest BizCustomer
	db.Table("biz_customer").Find(&biztest, "user_name", "플리test")
	fmt.Println(biztest.UserName, biztest.CompanyName)

	result := map[string]interface{}{}
	db.Table("biz_customer").Take(&result) // first, last는 모델과 매핑해야함
	fmt.Println(result)

	db.Table("biz_customer").Find(&bizcustom, "lower(user_name) Like ? AND created_at < ?", "%lower(KT)%", time.Now())
	//db.Table("biz_customer").Find(&bizcustom, "user_name", []string{"plea", "플리test"})
	//db.Table("biz_customer").Find(&bizcustom)  // 전체 데이터 다 가져오기
	for _, v := range bizcustom {
		fmt.Println(v)
	}

}
