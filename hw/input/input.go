package input

import (
	"bufio"
	"fmt"
	"os"
)

func Input() string {

	fmt.Println("데이터를 입력하세요")
	reader := bufio.NewReader(os.Stdin)
	inputData, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("잘못된 입력입니다.")
		Input()
	}
	return inputData

}
