package bplusmtree

import (
	"fmt"
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	testCount := 1000000
	bt := newBTree()
	start := time.Now()
	for i := testCount; i > 0; i-- {
		bt.Insert(i, "")
	}
	if bt.Count() != testCount {
		t.Error(bt.Count())
	}
	fmt.Println(time.Now().Sub(start))
	verifyTree(bt, testCount, t)
}

func TestSearch(t *testing.T) {
	testCount := 1000000
	bt := newBTree()

	for i := testCount; i > 0; i-- {
		bt.Insert(i, fmt.Sprintf("%d", i))
	}

	start := time.Now()
	for i := 1; i < testCount; i++ {
		v, ok := bt.Search(i)
		if !ok {
			t.Errorf("search: want = true, got = false")
		}
		if v != fmt.Sprintf("%d", i) {
			t.Errorf("search: want = %d, got = %s", i, v)
		}
	}
	fmt.Println(time.Now().Sub(start))
}

func verifyTree(b *BTree, count int, t *testing.T) {
	verifyRoot(b, t)

	for i := 0; i < b.Root().countNum(); i++ {
		verifyNode(b.Root().(*interiorNode).kcs[i].child, b.Root().(*interiorNode), t)
	}

	leftMost := findLeftMost(b.Root())

	if leftMost != b.First() {
		t.Errorf("bt.first: want = %p, got = %p", b.First(), leftMost)
	}

	verifyLeaf(leftMost, count, t)
}

// min child: 1
// max child: MaxKC
func verifyRoot(b *BTree, t *testing.T) {
	if b.Root().parent() != nil {
		t.Errorf("Root().parent: want = nil, got = %p", b.Root().parent())
	}

	if b.Root().countNum() < 1 {
		t.Errorf("Root().min.child: want >=1, got = %d", b.Root().countNum())
	}

	if b.Root().countNum() > MaxKC {
		t.Errorf("Root().max.child: want <= %d, got = %d", MaxKC, b.Root().countNum())
	}
}

func verifyNode(n node, parent *interiorNode, t *testing.T) {
	switch nn := n.(type) {
	case *interiorNode:
		if nn.countNum() < MaxKC/2 {
			t.Errorf("interior.min.child: want >= %d, got = %d", MaxKC/2, nn.countNum())
		}

		if nn.countNum() > MaxKC {
			t.Errorf("interior.max.child: want <= %d, got = %d", MaxKC, nn.countNum())
		}

		if nn.parent() != parent {
			t.Errorf("interior.parent: want = %p, got = %p", parent, nn.parent())
		}

		var last int
		for i := 0; i < nn.countNum(); i++ {
			key := nn.kcs[i].key
			if key != 0 && key < last {
				t.Errorf("interior.sort.key: want > %d, got = %d", last, key)
			}
			last = key

			if i == nn.countNum()-1 && key != 0 {
				t.Errorf("interior.last.key: want = 0, got = %d", key)
			}

			verifyNode(nn.kcs[i].child, nn, t)
		}

	case *leafNode:
		if nn.parent() != parent {
			t.Errorf("leaf.parent: want = %p, got = %p", parent, nn.parent())
		}

		if nn.countNum() < MaxKV/2 {
			t.Errorf("leaf.min.child: want >= %d, got = %d", MaxKV/2, nn.countNum())
		}

		if nn.countNum() > MaxKV {
			t.Errorf("leaf.max.child: want <= %d, got = %d", MaxKV, nn.countNum())
		}
	}
}

func verifyLeaf(leftMost *leafNode, count int, t *testing.T) {
	curr := leftMost
	last := 0
	c := 0

	for curr != nil {
		for i := 0; i < curr.countNum(); i++ {
			key := curr.kvs[i].key

			if key <= last {
				t.Errorf("leaf.sort.key: want > %d, got = %d", last, key)
			}
			last = key
			c++
		}
		curr = curr.next
	}

	if c != count {
		t.Errorf("leaf.count: want = %d, got = %d", count, c)
	}
}

func findLeftMost(n node) *leafNode {
	switch nn := n.(type) {
	case *interiorNode:
		return findLeftMost(nn.kcs[0].child)
	case *leafNode:
		return nn
	default:
		panic("")
	}
}
