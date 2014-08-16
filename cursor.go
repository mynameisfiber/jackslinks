// Package jackslinks is based off of http://arxiv.org/pdf/1408.1935v1.pdf
package jackslinks

import (
	"fmt"
)

// Cursor for list traversal
type Cursor struct {
	node *Node
}

// NewCursor will create a new cursor object at the given node
func NewCursor(node *Node) (*Cursor, error) {
	if node.IsHead() || node.IsTail() {
		return nil, ErrInvalidCursor
	}
	return &Cursor{node}, nil
}

// MoveToHead moves the cursor to the first element of the list
func (c *Cursor) MoveToHead() {
	for !c.node.prev.IsHead() {
		c.MoveLeft()
	}
}

// MoveToTail moves the cursor to the first element of the list
func (c *Cursor) MoveToTail() {
	for !c.node.next.IsTail() {
		c.MoveRight()
	}
}

// Print will show the LL associated with the current cursor to screen
func (c *Cursor) Print() {
	newC := &Cursor{c.node}
	fmt.Printf("CURSOR -> ")
	for {
		fmt.Printf("%v <-> ", newC.Get())
		if moved, _ := newC.MoveRight(); !moved {
			break
		}
	}
	fmt.Printf("TAIL\n")
}

// InsertAfter inserts a node after the current cursor
func (c *Cursor) InsertAfter(value interface{}) (bool, error) {
	newC := &Cursor{c.node.next}
	return newC.InsertBefore(value)
}

// InsertBefore inserts a node of given value before the cursor
func (c *Cursor) InsertBefore(value interface{}) (bool, error) {
	if c.node.IsHead() {
		return false, ErrInvalidCursor
	}
	for {
		y, yInfo, z, x, invDel, invIns := c.updateCursor()
		if invDel || invIns {
			return false, ErrInvalidCursor
		}
		zEffective := z
		if zEffective == nil {
			zEffective = newStateNode(TAIL)
		}
		nodes := [3]*Node{x, y, zEffective}
		oldInfo := [3]*Info{
			x.info,
			yInfo,
			zEffective.info,
		}
		if checkInfo(nodes, oldInfo) {
			dum := DummyInfo()
			newNode := &Node{&value, nil, x, nil, dum, ORDINARY}
			yCopy := &Node{y.Value, z, newNode, nil, dum, y.state}
			newNode.next = yCopy
			I := &Info{nodes, oldInfo, newNode, yCopy, false, INPROGRESS}
			if I.help() {
				c.node = yCopy
				return true, nil
			}
		}
	}
}

// Delete deletes the current node
func (c *Cursor) Delete() (bool, error) {
	if c.node.IsHead() || c.node.IsTail() {
		return false, nil
	}
	for {
		y, yInfo, z, x, invDel, _ := c.updateCursor()
		if invDel {
			return false, ErrInvalidCursor
		}
		nodes := [3]*Node{x, y, z}
		oldInfo := [3]*Info{x.info, yInfo, z.info}
		if checkInfo(nodes, oldInfo) {
			if y.IsTail() {
				return false, nil
			}
			I := &Info{nodes, oldInfo, z, x, true, INPROGRESS}
			if I.help() {
				if z.IsTail() {
					c.node = x
				} else {
					c.node = z
				}
				return true, nil
			}
		}
	}
}

// MoveLeft will move the cursor left
func (c *Cursor) MoveLeft() (bool, error) {
	y, _, _, x, invDel, _ := c.updateCursor()
	if invDel {
		return false, ErrInvalidCursor
	}
	if x.IsHead() {
		return false, nil
	}
	if x.state != ORDINARY && x.prev.next != x && x.next == y {
		if x.state == COPIED {
			c.node = x.newCopy
		} else {
			w := x.prev
			if w.IsHead() {
				return false, nil
			}
			c.node = w
		}
	} else {
		c.node = x
	}
	return true, nil
}

// MoveRight will move the cursor right
func (c *Cursor) MoveRight() (bool, error) {
	y, _, z, _, invDel, _ := c.updateCursor()
	if invDel {
		return false, ErrInvalidCursor
	}
	if y.next.IsTail() {
		return false, nil
	}
	c.node = z
	return true, nil
}

// Get the current value of the cursor
func (c *Cursor) Get() interface{} {
	y, _, _, _, invDel, _ := c.updateCursor()
	if invDel {
		return ErrInvalidCursor
	}
	return *(y.Value)
}

func (c *Cursor) updateCursor() (*Node, *Info, *Node, *Node, bool, bool) {
	invDel := false
	invIns := false
	for c.node.state != ORDINARY && c.node.prev.next != c.node {
		if c.node.state == COPIED {
			c.node = c.node.newCopy
			invIns = true
		} else {
			c.node = c.node.next
			invDel = true
		}
	}
	info := c.node.info
	return c.node, info, c.node.next, c.node.prev, invDel, invIns
}
