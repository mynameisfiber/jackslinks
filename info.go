package jackslinks

import (
	"sync/atomic"
	"unsafe"
)

// Info maintains state for nodes
type Info struct {
	nodes   [3]*Node
	oldInfo [3]*Info
	newNext *Node
	newPrev *Node
	remove  bool
	status  State
}

// DummyInfo creates a dummy info object
func DummyInfo() *Info {
	return &Info{
		nodes:   [3]*Node{nil, nil, nil},
		oldInfo: [3]*Info{nil, nil, nil},
		newNext: nil,
		newPrev: nil,
		remove:  false,
		status:  DUMMY,
	}
}

func checkInfo(nodes [3]*Node, oldInfo [3]*Info) bool {
	for _, info := range oldInfo {
		if info.status == INPROGRESS {
			info.help()
			return false
		}
	}
	for _, node := range nodes {
		if node.state != ORDINARY && node.state != HEAD && node.state != TAIL {
			return false
		}
	}
	for i := range nodes {
		if nodes[i].info != oldInfo[i] {
			return false
		}
	}
	return true
}

func (I *Info) help() bool {
	doPtrCAS := true
	for i := 0; i < 3 && doPtrCAS; i++ {
		atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&(I.nodes[i].info))),
			unsafe.Pointer(I.oldInfo[i]),
			unsafe.Pointer(I),
		)
		doPtrCAS = (I.nodes[i].info == I)
	}
	if doPtrCAS {
		if I.remove {
			I.nodes[1].state = MARKED
		} else {
			I.nodes[1].newCopy = I.newPrev
			I.nodes[1].state = COPIED
		}
		atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&(I.nodes[0].next))),
			unsafe.Pointer(I.nodes[1]),
			unsafe.Pointer(I.newNext),
		)
		atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&(I.nodes[2].prev))),
			unsafe.Pointer(I.nodes[1]),
			unsafe.Pointer(I.newPrev),
		)
		I.status = COMMITTED
	} else if I.status == INPROGRESS {
		I.status = ABORTED
	}
	return (I.status == COMMITTED)
}
