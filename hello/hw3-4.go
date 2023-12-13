package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type YamlConfig struct {
	Common struct {
		ServerName   string
		ServerNumber int
		Log          struct {
			FilePath    string
			Level       string
			ProcessName string
		}
		Mysql struct {
			Dsn     string
			MaxIdle int
			MaxConn int
		}
		MainBroker struct {
			Host string
		}
	}
	Server struct {
		Port string
	}
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
	var yaml YamlConfig

	var builder strings.Builder
	builder.WriteString("{ ")
	var upblankCnt int
	for fileScanner.Scan() {

		splitPoint := strings.IndexAny(fileScanner.Text(), ":")
		if splitPoint == -1 {
			panic("형식에 맞지 않습니다.")
		}
		splitKey := fileScanner.Text()[:splitPoint]
		splitVal := strings.TrimSpace(fileScanner.Text()[splitPoint+1:])

		blankCnt := strings.Count(strings.TrimRight(splitKey, " "), " ")

		if upblankCnt == blankCnt && blankCnt != 0 {
			builder.WriteString(", ")
		} else if upblankCnt != 0 && blankCnt == 0 {
			builder.WriteString("}},")
		} else if upblankCnt > blankCnt {
			builder.WriteString("}, ")
		}

		builder.WriteString("\"" + strings.TrimSpace(splitKey) + "\" : ")

		if splitVal == "" {
			builder.WriteString("{")
			upblankCnt = blankCnt
			continue
		}
		builder.WriteString(splitVal)
		upblankCnt = blankCnt
	}
	builder.WriteString("}}")
	fmt.Println(builder.String())

	err = json.Unmarshal([]byte(builder.String()), &yaml)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(yaml)
}
