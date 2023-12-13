package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type FileInfoHasTwoTypeValue struct {
	Key    string
	StrVal string
	IntVal int
}

type FileInfoHasTwoTypeValueTmp struct {
	Key string
	Val interface{}
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

	var fileInfoSlice []FileInfoHasTwoTypeValue
	tmp := make([]FileInfoHasTwoTypeValueTmp, 0)

	for fileScanner.Scan() {
		var fileInfo FileInfoHasTwoTypeValue
		var fileInfoTmp FileInfoHasTwoTypeValueTmp

		if !strings.Contains(fileScanner.Text(), ":") {
			panic("Key:Value 형식에 맞지 않습니다.")
		}

		splitStr := strings.Split(fileScanner.Text(), ":")

		fileInfo.Key = strings.TrimSpace(splitStr[0])
		fileInfoTmp.Key = strings.TrimSpace(splitStr[0])
		splitStr[1] = strings.TrimSpace(splitStr[1])

		if fileInfo.Key == "" {
			panic(fmt.Errorf("key 값이 없습니다"))
		} else if splitStr[1] == "" {
			panic(fmt.Errorf("%v의 Value 값이 없습니다", fileInfo.Key))
		}

		valueFirst := splitStr[1][0]
		valueEnd := splitStr[1][len(splitStr[1])-1]

		if string(valueFirst) == "\"" && string(valueEnd) == "\"" {
			splitStr[1] = strings.Trim(splitStr[1], "\"")
			fileInfo.StrVal = splitStr[1]
			fileInfoTmp.Val = splitStr[1]
		} else {
			strconvVal, err := strconv.Atoi(splitStr[1])
			if err != nil {
				panic("Value 값이 올바른 형식이 아닙니다.")
			} else {
				fileInfo.IntVal = strconvVal
				fileInfoTmp.Val = strconvVal
			}
		}

		fileInfoSlice = append(fileInfoSlice, fileInfo)
		tmp = append(tmp, fileInfoTmp)
	}

	for _, v := range fileInfoSlice {
		fmt.Printf("Key: %-20v, StrVal: %-20v, IntVal: %-20v\n", v.Key, v.StrVal, v.IntVal)
	}

	for _, v := range tmp {
		if reflect.TypeOf(v.Val).Kind() == reflect.String {
			fmt.Printf("Key: %s, Val:%s, Type:%s\n", v.Key, v.Val, reflect.TypeOf(v.Val))
		} else {
			fmt.Printf("Key: %s, Val:%d, Type:%s\n", v.Key, v.Val, reflect.TypeOf(v.Val))
		}
	}
}
