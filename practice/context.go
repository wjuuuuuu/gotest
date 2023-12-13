package main

import (
	"context"
	"fmt"
)

// context는 작업 명세서와 같은 역할로 작업 가능한 시간, 작업 취소 등 작업의 흐름을 제어
// ex) WithCancel, withDeadline, WithTimeout

// context.Background => 컨텍스트를 생성하는 방법
// 한 번 생성된 컨텍스트는 변경할 수 없다. 그래서 컨텍스트에 값을 추가할 때는 withValue 함수를 사용해서 새로운 컨텍스트 생성
// 컨텍스트의 값을 가져올 때는  context.Value(key)
// 취소 신호  context.witrhCancel(parent context)

type user struct {
	name string
}
type userKey int

func main() {
	u := user{name: "HOAN"}
	const uk userKey = 0

	ctx := context.WithValue(context.Background(), uk, &u) //상위 context(background), key, value(user의 주소값)

	if u, ok := ctx.Value(uk).(*user); ok {
		fmt.Println("User", u.name)
	}
}
