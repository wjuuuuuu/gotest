package main

import (
	"fmt"
	"hw/input"
	"os"
	"time"
)

func main() {
	now := time.Now()
	today := now.Format("20060102")
	thisTime := now.Format("150405")

	if err := os.Mkdir(today, os.ModePerm); err != nil {
		if os.IsExist(err) {
			fmt.Println("기존에 존재하는 디렉토리로 이동합니다.")
		} else {
			panic(err)
		}
	} else {
		fmt.Println(today, "새로운 디렉토리가 생성되었습니다.")
	}

	filePath := today + "/" + thisTime + ".txt"
	newFile, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	fmt.Println(thisTime + ".txt 파일이 생성되었습니다.")

	inputData := input.Input()

	_, err = newFile.WriteString(inputData)
	if err != nil {
		panic(err)
	}
	fmt.Println("작성이 완료되었습니다.")

	if err = newFile.Close(); err != nil {
		fmt.Println(err)
		return
	}

}
