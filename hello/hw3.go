package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	ServerName   string
	ServerNumber int
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

	filed := reflect.ValueOf(&Config{}).Elem()
	fiedMap := make(map[string]string)
	for i := 0; i < filed.NumField(); i++ {
		fiedMap[filed.Type().Field(i).Name] = filed.Type().Field(i).Type.String()
	}


	for fileScanner.Scan() {

		splitStr := strings.Split(fileScanner.Text(), ":")
		if len(splitStr) != 2 {
			panic("형식에 맞지 않습니다.")
		}

		splitStr[0] = strings.TrimSpace(splitStr[0])
		splitStr[1] = strings.TrimSpace(splitStr[1])

		//fmt.Println(splitStr[0], splitStr[1])

		// <--key 값 구조체 형식으로 변환 -->
		var strBuilder strings.Builder
		keyRunes := []rune(splitStr[0])
		upperChk := false
		for i, r := range keyRunes {
			if i == 0 || upperChk == true {
				strBuilder.WriteString(strings.ToUpper(string(r)))
				upperChk = false
			} else if string(r) == "_" {
				upperChk = true
			} else {
				strBuilder.WriteString(string(r))
			}
		}
		convKey := strBuilder.String()
		// <--key 값 구조체 형식으로 변환 끝 -->

		//value 값
		valueFirst := splitStr[1][0]
		valueEnd := splitStr[1][len(splitStr[1])-1]

		var configInstance Config

		if string(valueFirst) == "\"" && string(valueEnd) == "\"" {
			splitStr[1] = strings.Trim(splitStr[1], "\"")
		} else {
			strconvVal, err := strconv.Atoi(splitStr[1])
			if err != nil {
				panic("Value 값이 올바른 형식이 아닙니다.")
			} else {

				}
			}
		}




	}

}
