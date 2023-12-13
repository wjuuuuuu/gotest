package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type LevelConfig struct {
	Level int
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

	var Configs []LevelConfig
	var oldKey strings.Builder
	for fileScanner.Scan() {
		splitPoint := strings.IndexAny(fileScanner.Text(), ":")
		if splitPoint == -1 {
			fmt.Println(":가 없는 라인")
		}
		splitKey := fileScanner.Text()[:splitPoint]
		splitVal := fileScanner.Text()[splitPoint+1:]

		newBlankCount := strings.Count(splitKey, " ")
		level := newBlankCount/2 + 1

		if strings.TrimSpace(splitVal) == "" { // ex: common:  뒤에 아무것도 없을 때 oldKey에 key값을 넣고 그냥  return
			oldKey.WriteString(strings.TrimSpace(splitKey))
			continue
		}

		var config LevelConfig
		if oldKey.String() != "" { // oldkey가 있으면 oldkey를 key로 현재 key를 value로 config에 저장
			config.Level = level - 1
			config.Key = oldKey.String()
			config.Value = strings.TrimSpace(splitKey)
			Configs = append(Configs, config)
		}

		config.Level = level
		config.Key = strings.TrimSpace(splitKey)
		config.Value = strings.TrimSpace(splitVal)
		Configs = append(Configs, config)

		//splitStr[0] = strings.TrimSpace(splitStr[0])
		//splitStr[1] = strings.TrimSpace(splitStr[1])
		//
		//var strBuilder strings.Builder
		//keyRunes := []rune(splitStr[0])
		//upperChk := false
		//for i, r := range keyRunes {
		//	if i == 0 || upperChk == true {
		//		strBuilder.WriteString(strings.ToUpper(string(r)))
		//		upperChk = false
		//	} else if string(r) == "_" {
		//		upperChk = true
		//	} else {
		//		strBuilder.WriteString(string(r))
		//	}
		//}

	}
	sort.Slice(Configs, func(i, j int) bool {
		return Configs[i].Level < Configs[j].Level
	})
	for _, v := range Configs {
		fmt.Println(v.Level, v.Key, v.Value)
	}

}
