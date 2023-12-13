package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/v3/mem"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

type SwapMemory struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Used        uint64    `json:"used"`
	Free        uint64    `json:"free"`
	UsedPercent float64   `json:"used_percent"`
	CreateAt    time.Time `json:"create_at"`
}
type VirtualMemory struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Used        uint64    `json:"used"`
	Free        uint64    `json:"free"`
	UsedPercent float64   `json:"used_percent"`
	CreateAt    time.Time `json:"create_at"`
}

type Memory struct {
	SwapMemories    []SwapMemory    `json:"swap_memories"`
	VirtualMemories []VirtualMemory `json:"virtual_memories"`
}

func (SwapMemory) TableName() string {
	return "swap_memory"
}
func (VirtualMemory) TableName() string {
	return "virtual_memory"
}

func openDb() *gorm.DB {
	dsn := "root:./wjson./@tcp(localhost:3306)/sys?charset=utf8mb4&parseTime=True&loc=Local" // username:password/@tcp(host:port)/database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Db 연결에 실패하였습니다.")
	}
	return db
}

func saveData(db *gorm.DB) (err error) {
	tx := db.Begin()
	err = tx.AutoMigrate(&VirtualMemory{}, &SwapMemory{})
	if err != nil { // 테이블 생성 에러시 롤백
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func saveSwapData(db *gorm.DB, sMemories []SwapMemory) (err error) {
	tx := db.Begin()
	err = tx.Debug().Create(&sMemories).Error
	if err != nil {
		log.Println(err)
		tx.Rollback()
	}
	return tx.Commit().Error
}

func saveVirtualData(db *gorm.DB, vMemories []VirtualMemory) (err error) {
	tx := db.Begin()
	err = tx.Debug().Create(&vMemories).Error
	if err != nil {
		log.Println(err)
		tx.Rollback()
	}
	return tx.Commit().Error
}

func CreateDirectory() (string, error) {
	today := time.Now().Format("20060102")
	if err := os.Mkdir(today, os.ModePerm); err != nil {
		if !os.IsExist(err) {

			return "", err
		}
	}
	return today, nil
}
func CreateFile(directoryName string) (*os.File, error) {
	fileName := time.Now().Format("150405")
	newFile, err := os.Create(directoryName + "/memory_usage_" + fileName + ".txt")
	if err != nil {
		return nil, err
	}
	return newFile, nil
}

func WriteFile(newFile *os.File, sMemories []SwapMemory, vMemories []VirtualMemory) (err error) {
	_, err = newFile.WriteString("# SwapMemory Usage\n")
	if err != nil {
		return err
	}

	for _, s := range sMemories {
		_, err := fmt.Fprintf(newFile, "Time: %v, Used: %d, Free: %d, Used Percent: %f%%\n", s.CreateAt.Format("2006-01-02 15:04:05"), s.Used, s.Free, s.UsedPercent)
		if err != nil { // 작성에러
			return err
		}
	}

	_, err = newFile.WriteString("\n# VirtualMemory Usage\n")
	if err != nil {
		return err
	}

	for _, v := range vMemories {
		_, err := fmt.Fprintf(newFile, "Time: %s, Used: %d, Free: %d, Used Percent: %f%%\n", v.CreateAt.Format("2006-01-02 15:04:05"), v.Used, v.Free, v.UsedPercent)
		if err != nil { // 작성에러
			return err
		}
	}

	return nil
}

func findAllData(c echo.Context, db *gorm.DB) error {
	var sMemories []SwapMemory
	var vMemories []VirtualMemory

	errS := db.Debug().Find(&sMemories).Error
	errV := db.Debug().Find(&vMemories).Error

	if errS != nil || errV != nil {
		return c.String(http.StatusInternalServerError, "SERVER ERROR")
	}
	memory := Memory{
		sMemories,
		vMemories,
	}
	return c.JSON(http.StatusOK, memory)
}

func main() {

	sMemories := make([]SwapMemory, 0)
	vMemories := make([]VirtualMemory, 0)

	ch := make(chan string, 2)

	db := openDb()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return findAllData(c, db)
	})

	go func() {
		var metricError int

		tick := time.Tick(3 * time.Second)
		save := time.Tick(12 * time.Second)
		fin := time.After(26 * time.Second)

		for {
			select {
			case <-tick:
				if metricError > 10 {
					log.Println("Too Many Errors", metricError)
					//	wg.Done()
					return
				}
				v, err := mem.VirtualMemory()
				if err != nil {
					metricError++
					log.Println("Wrong data: ", v, err) // 에러난 데이터면 log 찍고 저장은 안함
					ch <- "v"
				} else {
					vMemories = append(vMemories, VirtualMemory{Used: v.Used, Free: v.Free, UsedPercent: v.UsedPercent, CreateAt: time.Now()}) // 정상데이터만 슬라이스에 저장
				}

				s, err := mem.SwapMemory()
				if err != nil {
					metricError++
					log.Println("Wrong data: ", s, err) // 에러난 데이터면 log 찍고 저장은 안함
					ch <- "s"
				} else {
					sMemories = append(sMemories, SwapMemory{Used: s.Used, Free: s.Free, UsedPercent: s.UsedPercent, CreateAt: time.Now()}) // 정상데이터만 슬라이스에 저장
				}

			case name := <-ch:
				time.Sleep(1 * time.Second)
				if name == "v" {
					v, err := mem.VirtualMemory()
					if err != nil {
						metricError++
						log.Println("Wrong data: ", v, err)
					} else {
						vMemories = append(vMemories, VirtualMemory{Used: v.Used, Free: v.Free, UsedPercent: v.UsedPercent, CreateAt: time.Now()}) // 정상데이터만 슬라이스에 저장
					}
				} else if name == "s" {
					s, err := mem.SwapMemory()
					if err != nil {
						metricError++
						log.Println("Wrong data: ", s, err)
					} else {
						sMemories = append(sMemories, SwapMemory{Used: s.Used, Free: s.Free, UsedPercent: s.UsedPercent, CreateAt: time.Now()}) // 정상데이터만 슬라이스에 저장
					}
				}

			case <-save:
				// db 저장
				err := saveData(db) // 테이블 생성시 에러면 return
				if err != nil {
					log.Println(err)
					return
				}
				err = saveSwapData(db, sMemories)
				if err != nil {
					log.Println(err)
				}
				err = saveVirtualData(db, vMemories)
				if err != nil {
					log.Println(err)
				}

				// 폴더, 파일 생성
				directoryName, err := CreateDirectory()
				if err != nil {
					directoryName, err = CreateDirectory() // 디렉토리 생성 에러시, 한 번 더 생성 -> 2차 에러시 그냥 진행
					fmt.Println(err)
				}

				newFile, err := CreateFile(directoryName)
				if err != nil {
					fileName := time.Now().Format("150405") // 파일 생성 에러시, 한 번 더 진행 -> 2차 에러시 return
					newFile, err = os.Create(directoryName + "/memory_usage_" + fileName + ".txt")
					if err != nil {
						fmt.Println(err)
						//		wg.Done()
						return
					}
				}
				// 데이터 출력
				err = WriteFile(newFile, sMemories, vMemories) // 파일입력 에러시 log 남기고 슬라이스 값 초기화 X
				if err != nil {
					fmt.Println(err)
				} else {
					// 메모리 비우기
					sMemories = make([]SwapMemory, 0) // 파일 입력 성공 시 슬라이스 값 초기화
					vMemories = make([]VirtualMemory, 0)
				}

			case <-fin:
				close(ch)
				return
			}

		}

	}()
	e.Logger.Fatal(e.Start(":1324"))

}
