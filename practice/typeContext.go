package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type UnmarshalTypeError struct {
	Value string       //Json Value
	Type  reflect.Type // 미리 선언할 수 없는 타입
}

// 모든 필드를 에러 메시지에서 사용하는가 검증한다. 그렇지 않다면 문제가 발생할 수 있다. 사용자 정의 에러 타입에 필드를 추가 해두어도 아래 메서드(Error)가 호출될 때 로그가 정상적으로 출력되지 않을 것
func (e *UnmarshalTypeError) Error() string {
	return "json : cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

type invalidUnmarshalError struct {
	Type reflect.Type
}

func (e *invalidUnmarshalError) Error() string { // unmarshal의 매개변수로 다른 타입이 들어왔을 때
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Ptr {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}

type User struct {
	Name int
}

func main() {
	var u User
	err := json.Unmarshal([]byte(`name":3}`), &u)
	if err != nil {
		switch e := err.(type) {
		case *UnmarshalTypeError:
			fmt.Printf("UnmarshalTypeError : Value[%s] Type[%v]\n", e.Value, e.Type)
		case *invalidUnmarshalError:
			fmt.Printf("InvalidUnmarshalError: Type[%v]\n", e.Type)
		default:
			fmt.Println(err)
		}
		return
	}
	fmt.Println("Name: ", u.Name)
}
