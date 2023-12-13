package benchmarking

import (
	"fmt"
	"testing"
)

var gs string

// Benchmark 테스트는 성능을 측정하는 기능 , Benchmark로 시작 ex) BenchmarkSum
// *testing.B 타입의 매개변수를 받는다

// Sprint 성능을 테스트
// 벤치마크 하려는 모든 코드는 b.N 반복문 안에 있어야한다.
// 처음 호출할 때 b.N은 1이고 테스트가 진행하면서 지속적으로 증가한다

// go test -run none -bench . -benchtime 3s -benchmem
func BenchmarkSprintBasic(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = fmt.Sprint("hello")
	}
	gs = s
}

func BenchmarkSprintfBasic(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = fmt.Sprintf("hello")
	}
	gs = s
}
