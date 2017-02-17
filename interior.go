package bplusmtree

import (
	"sort"
)
// Value 定义
type kc struct {
	key   int
	child node
}

// 预留一个空槽， 数组
type kcs [MaxKC + 1]kc

func (a *kcs) Len() int { return len(a) }

func (a *kcs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a *kcs) Less(i, j int) bool {
	// 处理数组中预留的空槽，value key 初始值是 0；应对 Search
	if a[i].key == 0 {
		return false
	}
	// 处理数组中预留的空槽，value key 初始值是 0；应对 Sort
	if a[j].key == 0 {
		return true
	}

	return a[i].key < a[j].key
}
// 中间节点数据结构定义
type interiorNode struct {
	kcs   kcs // 存储元素
	count int // 实际存储元素数目
	p     *interiorNode // 父亲节点
}

func newInteriorNode(p *interiorNode, largestChild node) *interiorNode {
	i := &interiorNode{
		p:     p,
		count: 1,
	}

	if largestChild != nil {
		i.kcs[0].child = largestChild
	}
	return i
}



// 从该中间节点找到 key 元素应该存储的位置
func (in *interiorNode) find(key int) (int, bool) {
	// 定义查询方法，这里只需要 ">"
	c := func(i int) bool { return in.kcs[i].key > key }
	// 查询
	i := sort.Search(in.count-1, c)

	return i, true
}
// 判断是否达到中间节点的最大元素数目限制 MaxKC
func (in *interiorNode) full() bool { return in.count == MaxKC }
// 返回中间节点的父亲节点
func (in *interiorNode) parent() *interiorNode { return in.p }
// 设置中间节点的父亲节点
func (in *interiorNode) setParent(p *interiorNode) { in.p = p }

func (in *interiorNode) countNum() int { return in.count }


// 插入中间节点
func (in *interiorNode) insert(key int, child node) (int, *interiorNode, bool) {
	// 确定 key 在该中间节点应该存储的位置
	i, _ := in.find(key)
	// 中间节点没有达到数量限制
	if !in.full() {
		// 将元素插入中间节点
		copy(in.kcs[i+1:], in.kcs[i:in.count])
		// 设置子节点分裂后产生的元素 为当前位置 i 的key
		in.kcs[i].key = key
		// 设置子节点以及子节点设置父亲节点
		in.kcs[i].child = child
		child.setParent(in)
		// 元素计数加一
		in.count++
		return 0, nil, false
	}

	// 达到数量限制，则在最右侧保留的空槽追加该节点
	in.kcs[MaxKC].key = key
	in.kcs[MaxKC].child = child
	// 子节点设置父亲节点
	child.setParent(in)
	// 中间节点分裂
	next, midKey := in.split()

	return midKey, next, true
}

func (in *interiorNode) split() (*interiorNode, int) {
	// 节点排序，把新插入的节点防到正确的位置
	sort.Sort(&in.kcs)

	// 获取中间元素的位置，并设置 Value 的 子节点和 key
	midIndex := MaxKC / 2
	midChild := in.kcs[midIndex].child
	midKey := in.kcs[midIndex].key

	// 创建一个新没有父亲节点(第一个 nil)的中间节点
	next := newInteriorNode(nil, nil)
	// 将中间元素的右侧数组拷贝到新的分裂节点
	copy(next.kcs[0:], in.kcs[midIndex+1:])
	// 设置分裂节点的 Count
	next.count = MaxKC - midIndex
	// 更新分裂节点中所有元素子节点的父亲节点
	for i := 0; i < next.count; i++ {
		next.kcs[i].child.setParent(next)
	}

	// 更新原始节点的参数，将中间元素放进原始节点
	in.count = midIndex + 1
	in.kcs[in.count-1].key = 0
	in.kcs[in.count-1].child = midChild
	midChild.setParent(in)
	// 返回分裂后产生的中间节点和中间元素的 key，供父亲节点插入
	return next, midKey
}
