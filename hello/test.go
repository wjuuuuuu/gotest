package main

import (
	"fmt"
	"strconv"
	"strings"
)

// 문자열이 숫자로만 이루어져 있는지 확인
func IsNum(inputNum string) bool {
	_, err := strconv.Atoi(inputNum)
	if err != nil {
		return false
	}
	return true
}

// 앞에 0이 있으면 다 떼고 하나의 0을 붙여줌
func FirstOnlyOneZero(inputNum string) (string, bool) {
	if !IsNum(inputNum) {
		return "숫자가 아닙니다.", false
	}
	runes := []rune(inputNum)
	var zeroCount int
	for _, r := range runes {
		if r != '0' {
			break
		} else {
			zeroCount++
		}
	}
	inputNum = "0" + inputNum[zeroCount:]
	return inputNum, true
}

// 국번 값이 0으로 시작하면 0은 삭제 처리
func FinalPhoneNum(inputNum string, chk bool) string {
	if !chk {
		return inputNum
	}
	inputNumSlice := strings.Split(inputNum, "")
	second := inputNumSlice[1]
	dial := ""
	fullLen := len(inputNum)

	switch second {
	case "2":
		fullLen = fullLen - 2
		dial = inputNumSlice[2]
	case "3", "4", "5", "6", "7":
		fullLen = fullLen - 3
		dial = inputNumSlice[3]
	default:
		return "잘못된 지역번호 입니다."
	}

	if dial == "0" {
		if second == "2" {
			// 복잡, 함수로 만들기(코드 중복)
			inputNumSlice = inputNumSlice[:2+copy(inputNumSlice[2:], inputNumSlice[3:])]
			fullLen = fullLen - 1
		} else {
			inputNumSlice = inputNumSlice[:3+copy(inputNumSlice[3:], inputNumSlice[4:])]
			fullLen = fullLen - 1
		}
	}

	// 자리수가 7인 경우와 8인 경우 지역번호, 국번, 전화번호 나눠서 출력
	inputNum = strings.Join(inputNumSlice, "")
	if fullLen == 7 {
		// 복잡하고 함수로 만들어서..
		fmt.Println("지역번호: ", inputNum[:len(inputNum)-7], "국번: ", inputNum[len(inputNum)-7:len(inputNum)-4], "전화번호: ", inputNum[len(inputNum)-4:])
	} else if fullLen == 8 {
		fmt.Println("지역번호: ", inputNum[:len(inputNum)-8], "국번: ", inputNum[len(inputNum)-8:len(inputNum)-4], "전화번호: ", inputNum[len(inputNum)-4:])
	} else {
		return "전화번호 길이가 다릅니다."
	}
	return inputNum

}

func main() {
	var inputNum string = "0700123987678"
	// 지역번호 추출, 국번 추출, 전화번호 추출 3개의 함수.

	inputNum, chk := FirstOnlyOneZero(inputNum)
	inputNum = FinalPhoneNum(inputNum, chk)
	fmt.Println(inputNum)
}
