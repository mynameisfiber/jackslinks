package jackslinks

import (
	"errors"
)

// State flag for a node
type State uint8

// Possible states that a node can be in
const (
	DUMMY = iota
	COPIED
	MARKED
	ORDINARY
	INPROGRESS
	COMMITTED
	ABORTED
	HEAD
	TAIL
)

var (
	// ErrInvalidCursor indicates an operation is attempted on an invalid cursor state
	ErrInvalidCursor = errors.New("Invalid cursor")
)

// Node instance
type Node struct {
	Value   *interface{}
	next    *Node
	prev    *Node
	newCopy *Node
	info    *Info
	state   State
}

// NewNode creates a new node
func NewNode(value interface{}) *Node {
	node := &Node{
		Value:   &value,
		next:    newStateNode(TAIL),
		prev:    newStateNode(HEAD),
		newCopy: nil,
		info:    DummyInfo(),
		state:   ORDINARY,
	}
	node.next.prev = node
	node.prev.next = node
	return node
}

func newStateNode(state State) *Node {
	node := Node{
		Value:   nil,
		next:    nil,
		prev:    nil,
		newCopy: nil,
		info:    DummyInfo(),
		state:   state,
	}
	return &node
}

// IsHead checks if current node is head
func (n *Node) IsHead() bool {
	return n.state == HEAD
}

// IsTail checks if current node is head
func (n *Node) IsTail() bool {
	return n.state == TAIL
}
