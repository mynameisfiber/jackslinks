package jackslinks

import (
	"sync"
	"testing"
)

func TestDeleteTop(t *testing.T) {
	node := NewNode(1)
	cursor, _ := NewCursor(node)
	cursor.InsertAfter(2)
	cursor.Delete()
	if cursor.Get().(int) != 2 {
		t.Errorf("Didn't move to the right place")
	}
}

func TestDeleteBottom(t *testing.T) {
	node := NewNode(1)
	cursor, _ := NewCursor(node)
	cursor.InsertBefore(2)
	cursor.Delete()
	if cursor.Get().(int) != 2 {
		t.Errorf("Didn't move to the right place")
	}
}

func TestDLL(*testing.T) {
	node := NewNode("test1")
	cursor, _ := NewCursor(node)
	cursor.InsertBefore("test2")
	cursor.InsertBefore("test3")
	cursor.MoveToHead()
	cursor.InsertBefore("test4")
	cursor.InsertAfter("test5")
	cursor.Print()

	cursor.MoveToTail()
	cursor.InsertAfter("test6")
	cursor.Print()

	cursor.Delete()
	cursor.Print()

	cursor.MoveToTail()
	cursor.Delete()
	cursor.MoveToHead()
	cursor.Print()

	cursor.MoveToTail()
	cursor.InsertAfter("test7")
	cursor.MoveToHead()
	cursor.Print()
}

func TestHolistic(*testing.T) {
	workers := 5
	itemsPerWorker := 50
	var waitGroup sync.WaitGroup

	rootNode := NewNode(-1)
	list, _ := NewCursor(rootNode)
	for i := 0; i < workers; i++ {
		waitGroup.Add(1)
		go func(id int) {
			for j := id * itemsPerWorker; j < (id+1)*itemsPerWorker; j++ {
				list.InsertBefore(j)
				if j%(id+3) == 0 {
					list.Delete()
				}
			}
			waitGroup.Done()
		}(i)
	}
	waitGroup.Wait()

	list.MoveToHead()
	list.Print()
}
