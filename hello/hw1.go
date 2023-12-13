package main

// 23/10/16
// 변수명, 메서드명 길어도 상관없으니 누가봐도 알아볼수 있도록 고민해서 짓기
//	작성한 코드에 밑줄이 생기면 뭔가 더 나은 방법이 있는지, 발생될 문제가 있는지 고민해보기

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TestStructure struct {
	Key   string
	Value string
}

func main() {
	file, err := os.Open("C:\\Users\\wjuuu\\Desktop\\key_value_parsing.txt")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	var testSlice []TestStructure
	for fileScanner.Scan() {
		var test TestStructure

		if !strings.Contains(fileScanner.Text(), ":") {
			panic("Key:Value 형식에 맞지 않습니다.")
		}

		splitStr := strings.Split(fileScanner.Text(), ":")

		test.Key = strings.Trim(splitStr[0], " ")
		test.Value = strings.Trim(splitStr[1], " ")

		if test.Key == "" {
			panic(fmt.Errorf("key 값이 없습니다"))
		} else if test.Value == "" {
			panic(fmt.Errorf("%v의 Value 값이 없습니다", test.Key))
		}

		testSlice = append(testSlice, test)
	}

	for _, v := range testSlice {
		fmt.Printf("Key : %-15v, Val: %-15v\n", v.Key, v.Value)
	}

}
