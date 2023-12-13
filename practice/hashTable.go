package main

import (
	"fmt"
	"time"
)

// 해시테이블
// 충돌의 경우 chainig을 array로 수행

var bucket_size = 10

type node struct {
	key   int
	value int
	next  *node // 다음 노드
}

type hashTable struct {
	bucket map[int][]node
}

func (n *node) createNode(key, value int) node {
	n.key = key
	n.value = value
	n.next = nil

	return *n
}
func (n *node) putNext(node *node) {
	n.next = node
}

func hashFunction(key int) int {
	return key % bucket_size
}

func (h *hashTable) createBucket() hashTable {
	h.bucket = make(map[int][]node, 1)
	return *h
}

func (h *hashTable) addTable(key, value int) {
	hashIndex := hashFunction(key)
	var newNode node
	newNode.createNode(key, value)

	slotCount := len(h.bucket[hashIndex])

	if slotCount == 0 {
		h.bucket[hashIndex] = append(h.bucket[hashIndex], newNode)
	} else {
		var lastNode *node
		for i := range h.bucket[hashIndex] {
			if h.bucket[hashIndex][i].next == nil {
				lastNode = &h.bucket[hashIndex][i]
				break
			}
		}
		if lastNode != nil {
			lastNode.putNext(&newNode)
		}
		h.bucket[hashIndex] = append(h.bucket[hashIndex], newNode)

	}
}

func (h hashTable) search(key int) {
	hashIndex := hashFunction(key)
	var nodes []node = h.bucket[hashIndex]
	var findNode node
	for _, v := range nodes {
		if v.key == key {
			findNode = v
			break
		}
	}
	if &findNode != nil {
		fmt.Printf("키는 [%d], 값은 [%d] 입니다", findNode.key, findNode.value)
	} else {
		fmt.Println("일치하는 키가 없습니다.")
	}

}

func main() {
	var hashTable1 hashTable
	hashTable1.createBucket()

	hashTable1.addTable(111, 30)
	hashTable1.addTable(121, 40)
	hashTable1.addTable(131, 50)
	hashTable1.addTable(112, 60)
	hashTable1.addTable(122, 70)
	hashTable1.addTable(123, 80)
	hashTable1.addTable(131, 90)
	hashTable1.addTable(132, 100)
	hashTable1.addTable(133, 110)

	start := time.Now()
	hashTable1.search(132)
	end := time.Now().Sub(start)
	fmt.Println(end)

}
