package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type configSt struct {
	common struct {
		serverName   string `tag:"ServerName"`
		serverNumber int    `tag:"ServerNumber"`
	}
}

func main() {

	obj := &configSt{}
	s := reflect.ValueOf(obj).Elem()
	typeOfT := s.Type() // main.configSt
	//	key := "ServerName"
	for i := 0; i < s.NumField(); i++ {
		a := typeOfT.Field(i).Type // struct

		fmt.Println("###")
		fmt.Println(a.Field(i).Tag.Lookup("tag"))
		fmt.Println(a.Field(i).Tag.Get("tag"))

		fmt.Println(a, s.Field(i).Kind())
		if s.Field(i).Kind() == reflect.Struct {

			for k := 0; k < a.NumField(); k++ {
				b := a.Field(k).Type // string, int
				fmt.Println(b, a.Field(k).Tag, b.Kind())

			}
		}
	}

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

	var upblankCnt int
	structures := make([]string, 0)
	for fileScanner.Scan() {

		splitPoint := strings.IndexAny(fileScanner.Text(), ":")
		if splitPoint == -1 {
			panic("형식에 맞지 않습니다.")
		}
		splitKey := fileScanner.Text()[:splitPoint]
		splitVal := strings.TrimSpace(fileScanner.Text()[splitPoint+1:])

		blankCnt := strings.Count(strings.TrimRight(splitKey, " "), " ")

		fmt.Println(splitKey, splitVal, blankCnt, upblankCnt)

		if splitVal == "" {
			upblankCnt = blankCnt
			structures = append(structures, strings.TrimSpace(splitKey))
			continue
		}

		if upblankCnt < blankCnt {
			makeStruct := structures[len(structures)-1]
			fmt.Println(makeStruct, "만들어야할 상위 구조체")
		}

	}
}
