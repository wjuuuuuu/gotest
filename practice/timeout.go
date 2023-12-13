package main

// 프로그램 모니터링, 너무 오래돌면 종료
import (
	"errors"
	"log"
	"os"
	"os/signal"
	"time"
)

const timeoutSeconds = 3 * time.Second

var (
	sigChan  = make(chan os.Signal, 1)    // 운영체제 신호를 받는 버퍼 1인 채널
	timeout  = time.After(timeoutSeconds) // 시간 초과 후
	complete = make(chan error)           // 처리가 끝난
	shutdown = make(chan struct{})        // 시스템 전체 알림
)

func main() {
	log.Println("Stating Process")
	signal.Notify(sigChan, os.Interrupt) // sigChan 채널에게 os.Interrupt와 관련있는 신호가 보이면 데이터 신호를 보내라.
	// 신호를 받지 못하면 기다리지 않기 때문에 버퍼가 있는 채널을 사용해서 1개의 신호를 보장받는다

	log.Println("Launching Processors")
	go processor(complete)

ControlLoop:
	for {
		select {
		case <-sigChan:
			log.Println("OS INTERRUPT")
			close(shutdown) // interrupt 관련 신호가 오면 채널 닫기
			sigChan = nil   // 닫긴 채널에서 재호출시 panic에 빠지기 때문에 nil 값으로 만들어줌
		case <-timeout:
			log.Println("Timeout - Killing program")
			os.Exit(1)
		case err := <-complete:
			log.Printf("Task Completed Error[%s]", err)
			break ControlLoop
		}
	}
	log.Println("Process Ended")
}

func processor(complete chan<- error) { // 신호 보내기 전용
	log.Println("Processor - Starting")
	var err error
	defer func() {
		if r := recover(); r != nil {
			log.Println("Processor - Panic", r)
		}
		complete <- err
	}()
	err = doWork()
	log.Println("Processor - Completed")
}

func doWork() error {
	log.Println("Processor - Task 1")
	time.Sleep(2 * time.Second)

	if checkShutdown() {
		return errors.New("EARLY SHUT DOWN")
	}

	log.Println("Processor - Task 2")
	time.Sleep(1 * time.Second)
	if checkShutdown() {
		return errors.New("EARLY SHUT DOWN")
	}

	log.Println("Processor - Task 3")
	time.Sleep(1 * time.Second)
	return nil
}

func checkShutdown() bool {
	select {
	case <-shutdown:
		log.Println("checkShutdown - Shutdown Early")
		return true
	default:
		return false
	}
}
