package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type result struct {
	id  int
	op  string
	err error
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	const routines = 10
	const inserts = routines * 2

	ch := make(chan result, inserts) // 버퍼가 20개인 채널

	waitInserts := inserts // waitGroup 대신 사용

	for i := 0; i < routines; i++ {
		go func(id int) {
			ch <- insertUser(id)  // db USER table에 insert 하는 작업
			ch <- insertTrans(id) // db TRANS table에 insert 하는 작업
		}(i)
	}

	for waitInserts > 0 {
		r := <-ch
		log.Printf("N %d ID: %d OP: %s ERR: %v", waitInserts, r.id, r.op, r.err)
		waitInserts--
	}
	log.Println("Inserts Complete")
}

func insertUser(id int) result {
	r := result{
		id: id,
		op: fmt.Sprintf("insert USERS value (%d)", id),
	}
	if rand.Intn(10) == 0 {
		r.err = fmt.Errorf("Unable to insert %d into USER table", id)
	}
	return r
}

func insertTrans(id int) result {
	r := result{
		id: id,
		op: fmt.Sprintf("insert TRANS value (%d)", id),
	}
	if rand.Intn(10) == 0 {
		r.err = fmt.Errorf("Unable to insert %d into TRANS table", id)
	}
	return r
}
