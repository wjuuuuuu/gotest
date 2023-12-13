package main

import (
	"fmt"
	"time"
)

type Node struct {
	key   int
	value int
	left  *Node
	right *Node
}
type Bst struct {
	root *Node
}

func (bst *Bst) Search(key int) *Node {
	return search(bst.root, key)

}

func search(n *Node, k int) *Node {
	if n == nil || n.key == k {
		return n
	}
	if k < n.key {
		return search(n.left, k)
	}
	return search(n.right, k)
}
func (bst *Bst) Insert(key int, value int) {
	n := &Node{key, value, nil, nil}
	if bst.root == nil {
		bst.root = n
	} else {
		insert(bst.root, n)
	}
}
func insert(n *Node, new *Node) {
	if new.key < n.key {
		if n.left == nil {
			n.left = new
		} else {
			insert(n.left, new)
		}
	} else {
		if n.right == nil {
			n.right = new
		} else {
			insert(n.right, new)
		}
	}
}

func main() {
	var binarySerchTree Bst
	binarySerchTree.Insert(111, 30)
	binarySerchTree.Insert(121, 40)
	binarySerchTree.Insert(131, 50)
	binarySerchTree.Insert(112, 60)
	binarySerchTree.Insert(122, 70)
	binarySerchTree.Insert(123, 80)
	binarySerchTree.Insert(131, 90)
	binarySerchTree.Insert(132, 100)
	binarySerchTree.Insert(133, 110)

	start := time.Now()
	v := binarySerchTree.Search(132)
	end := time.Now().Sub(start)
	fmt.Println(v.key, v.value, end)

}
