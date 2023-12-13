package main

// 에러처리 확실하게 하기, 기능별로 메서드로 묶어서 빼기, log 남길때는 명확하게 남기기, 중복되는 변수명 사용 안하기
// panic 쓰지말고 error를 return
import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type yamlConfig struct {
	Common Common `yaml:"Common"`
	Server Server `yaml:"Server"`
}
type Common struct {
	ServerName   string     `yaml:"ServerName"`
	ServerNumber int        `yaml:"ServerNumber"`
	Log          Log        `yaml:"Log"`
	Mysql        Mysql      `yaml:"Mysql"`
	MainBroker   MainBroker `yaml:"MainBroker"`
}

type Server struct {
	Port         string `yaml:"Port"`
	ServerName   string `yaml:"ServerName"`
	ServerNumber int    `yaml:"ServerNumber"`
	Host         string `yaml:"Host"`
}

type Log struct {
	FilePath    string `yaml:"FilePath"`
	Level       string `yaml:"Level"`
	ProcessName string `yaml:"ProcessName"`
}
type Mysql struct {
	Dsn     string `yaml:"Dsn"`
	MaxIdle int    `yaml:"MaxIdle"`
	MaxConn int    `yaml:"MaxConn"`
}
type MainBroker struct {
	Host string `yaml:"Host"`
}

func main() {
	var config yamlConfig
	file, err := openFile("C:\\Users\\wjuuu\\Desktop\\key_value_parsing.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	convertFileToMap := readFileAndParsToMap(file)

	err = convertMapToStruct(convertFileToMap, &config, "")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v", config)

	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}

}

func openFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func readFileAndParsToMap(file *os.File) map[string]interface{} {
	var lineNum, upBlankCnt int
	var upKey string
	convertStrToMap := make(map[string]interface{}, 0)

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		lineNum++

		if lineBlank := strings.TrimSpace(fileScanner.Text()) == ""; lineBlank {
			continue
		}

		if strings.HasPrefix(fileScanner.Text(), "\t") == false {
			// 최상위
			// 상태값 초기화
		} else {
			// 하위
			// 1 / 2/ 3
			// 현재 진행중인 상위 구조체 정보 저장 등...
		}

		splitPoint := strings.IndexAny(fileScanner.Text(), ":")
		if splitPoint == -1 {
			log.Panicf("형식에 맞지 않습니다. line: %d, Content: %s", lineNum, fileScanner.Text())
		}
		splitKey := fileScanner.Text()[:splitPoint]
		splitVal := strings.TrimSpace(fileScanner.Text()[splitPoint+1:])

		blankCnt := strings.Count(strings.TrimRight(splitKey, " "), " ")

		if splitVal == "" {
			if blankCnt == 0 {
				upKey = strings.TrimSpace(splitKey)
			} else if upBlankCnt >= blankCnt {
				lastIndex := strings.LastIndex(upKey, ".")
				upKey = upKey[:lastIndex+1] + strings.TrimSpace(splitKey)
			} else {
				upKey += "." + strings.TrimSpace(splitKey)
			}
			upBlankCnt = blankCnt
			continue
		}

		splitKey = strings.TrimSpace(splitKey)
		splitVal = strings.TrimSpace(strings.Trim(splitVal, "\""))

		convertStrToMap[upKey+"."+splitKey] = splitVal
	}
	return convertStrToMap
}

func convertMapToStruct(m map[string]interface{}, structure interface{}, upStruct string) error {
	structValue := reflect.ValueOf(structure).Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)
		fieldTag := fieldType.Tag.Get("yaml")

		switch fieldType.Type.Kind() {
		case reflect.Struct:
			var sendStr string
			if upStruct == "" {
				sendStr = fieldType.Name + "."
			} else {
				sendStr = upStruct + fieldType.Name + "."
			}

			if err := convertMapToStruct(m, field.Addr().Interface(), sendStr); err != nil {
				return err
			}
		case reflect.String:
			if val, ok := m[upStruct+fieldTag]; ok {
				field.Set(reflect.ValueOf(val))
			}
		case reflect.Int:
			if val, ok := m[upStruct+fieldTag]; ok {
				intVal, err := strconv.Atoi(val.(string))
				if err != nil {
					return err
				}
				field.Set(reflect.ValueOf(intVal))
			}
		}
	}
	return nil
}
