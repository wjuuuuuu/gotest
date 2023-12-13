package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type BizCustomer struct {
	gorm.Model
	UserInfo        string
	UserName        string
	CompanyName     string
	CustomerGroupId string
}

type JoinResult struct {
	CustomerGroupId   string
	CustomerGroupName string
	Token             string
}

func main() {
	dns := "dxm:Cjswo.123@tcp(im.plea.kr:13306)/dxm?parseTime=True"
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic("db connection fail")
	}

	var bizcustom []BizCustomer

	db.Table("biz_customer").Where("user_name Like ?", "%kt%").Find(&bizcustom)

	db.Table("biz_customer").Not("user_name Like ?", "%kt%").Find(&bizcustom)
	db.Table("biz_customer").Not("id", []int{1, 5, 13, 15, 16, 18, 19, 20}).Find(&bizcustom)

	db.Table("biz_customer").Where("user_name = ?", "mono").Or("company_name= ?", "KT").Find(&bizcustom)

	db.Table("biz_customer").Select("COALESCE(customer_group_id,company_name) as user_name").Find(&bizcustom)
	// customer_group_id 컬럼의 값이 null이면 company_name 컬럼의 값을 반환하고 user_name으로 대입하여 출력

	db.Table("biz_customer").Where("user_name Like ?", "%test%").Limit(10).Offset(20).Order("id desc").Find(&bizcustom)
	// 정렬하고 20번째 부터 10개의 결과값을 출력한다. => 10개씩 나눠진 데이터의 2페이지

	db.Table("biz_customer").Where("user_name Like ?", "%test%").Group("DATE(created_at)").Order("id desc").Find(&bizcustom)

	for _, v := range bizcustom {
		fmt.Println(v.ID, v.UserName, v.UserInfo, v.CompanyName)
	}

	fmt.Println()

	// iteration , rows에 가져온 데이터들을 담고 후에 반복문을 통해 하나씩 구조체에 담든지.. 출력을 하든지..
	rows, err := db.Table("biz_customer_group bcg").Distinct("bcg.customer_group_id", "bcg.customer_group_name", "IFNULL(bcg.token, 'HAS NO TOKEN') as token").
		Joins("inner join biz_customer bc on bcg.customer_group_id = bc.customer_group_id").Rows()

	for rows.Next() {
		var joinresult JoinResult
		db.ScanRows(rows, &joinresult)
		fmt.Println(joinresult)
	}

}
