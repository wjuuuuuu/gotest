package main

// 기준이 되는 시간이 하나인 방식으로 하는게 좋다
// ex) tick을 기준으로 한다면 printDataToFile은 굳이 고루틴으로 만들지 않아도 될듯 하다..
import (
	"fmt"
	"github.com/mackerelio/go-osstat/memory"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	ch := make(chan string)
	wg.Add(2)
	go makeSendData(ch)
	go printDataToFile(ch)
	fmt.Println(time.Now(), "Main")
	wg.Wait()
}

func makeSendData(ch chan string) {
	checkMemory := time.Tick(1 * time.Second)
	sendMemory := time.Tick(5 * time.Second)
	terminate := time.After(21 * time.Second)
	builder := strings.Builder{}

	for {
		select {
		case <-checkMemory:
			memoryInfo, err := memory.Get()
			if err != nil {
				fmt.Println(err)
				return
			}
			builder.WriteString(time.Now().Format("2006-01-02 15:04:05") + "-> " +
				strconv.FormatUint(memoryInfo.Used, 10) + "\n")

		case <-sendMemory:
			ch <- builder.String()
			builder.Reset()

		case <-terminate:
			wg.Done()
			close(ch)
			return
		}
	}
}

func printDataToFile(ch chan string) {
	now := time.Now()
	today := now.Format("20060102")
	if err := os.Mkdir(today, os.ModePerm); err != nil {
		if os.IsExist(err) {
			fmt.Println("경로가 이미 있습니다")
		} else {
			panic(err)
		}
	}

	for memoryInfo := range ch {
		thisTime := time.Now().Format("150405")
		newFile, err := os.Create(today + "/" + thisTime + ".txt")
		if err != nil {
			panic(err)
		}

		_, err = newFile.WriteString(memoryInfo)
		if err != nil {
			panic(err)
		}
		err = newFile.Close()
		if err != nil {
			panic(err)
		}
	}
	wg.Done()
}
