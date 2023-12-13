package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var (
	data      []string
	rmMutex   sync.RWMutex
	readCount int64
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func writer(i int) {
	rmMutex.Lock() // 쓰기 뮤텍스 잠금
	{
		rc := atomic.LoadInt64(&readCount) // 현재 readCount 캡쳐 , 항상 0이여야함. (잠금이니까)
		fmt.Printf("***> : Performing Write : RCount[%d]\n", rc)
		data = append(data, fmt.Sprintf("String: %d", i)) // 슬라이스에 새 문자열 추가
	}
	rmMutex.Unlock()
}
func reader(id int) {
	rmMutex.RLock() // 읽기 뮤텍스 잠금
	{
		rc := atomic.AddInt64(&readCount, 1) // readCount 1 증가

		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond) // 잠시 휴식 -> 이때 다른 고루틴이 실행 rc 증가
		//time.Sleep(time.Duration(10 * time.Millisecond))
		fmt.Printf("%d : Performing Read : Length[%d] RCount[%d]\n", id, len(data), rc)
		atomic.AddInt64(&readCount, -1) // readCount 1 감소
	}
	rmMutex.RUnlock()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for i := 0; i < 10; i++ { // 10개의 쓰기
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			writer(i)
		}
		wg.Done()
	}()

	for i := 0; i < 5; i++ {
		go func(i int) {
			for { // 영원한 5개의 읽기
				reader(i)
			}
		}(i)
	}
	wg.Wait()
	fmt.Println("Program Complete")
}
