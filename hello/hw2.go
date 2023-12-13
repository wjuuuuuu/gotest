package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*
Name: "Chaeun"
Email: "codjs@plea.kr"
Comapany: "Plea"
Number: 19

파일을 읽고 구조체에 담아서 출력하기2
- string은 "" 안에 있는 데이터로 구분함
- 어떤 구조체에 담을지 고민하기
*/

type FileInfo struct {
	Key   string
	Value interface{}
}

func main() {
	file, err := os.Open("C:\\Users\\wjuuu\\Desktop\\key_value_parsing.txt")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var fileInfoSlice []FileInfo

	for fileScanner.Scan() {
		var fileInfo FileInfo

		if !strings.Contains(fileScanner.Text(), ":") {
			panic("Key:Value 형식에 맞지 않습니다.")
		}

		splitStr := strings.Split(fileScanner.Text(), ":")

		fileInfo.Key = strings.TrimSpace(splitStr[0])
		fileInfo.Value = strings.TrimSpace(splitStr[1])

		if fileInfo.Key == "" {
			panic(fmt.Errorf("key 값이 없습니다"))
		} else if fileInfo.Value == "" {
			panic(fmt.Errorf("%v의 Value 값이 없습니다", fileInfo.Key))
		}

		fileInfoSlice = append(fileInfoSlice, fileInfo)
	}

	for _, v := range fileInfoSlice {
		fmt.Printf("Key: %v, Val: %v\n", v.Key, v.Value)
	}

}
