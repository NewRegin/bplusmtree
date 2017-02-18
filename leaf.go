package bplusmtree

import (
	"sort"
)

// Value 定义
type kv struct {
	key   int
	value string
}

// 叶子节点的存储数组，[M]value
type kvs [MaxKV]kv

func (a *kvs) Len() int           { return len(a) }
func (a *kvs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a *kvs) Less(i, j int) bool { return a[i].key < a[j].key }

// 叶子节点的数据结构
type leafNode struct {
	kvs   kvs           // 存储元素
	count int           // 实际存储元素数目
	next  *leafNode     // 右边第一个叶子节点（右指针）
	p     *interiorNode // 父亲节点，中间节点
}

// 创建新的叶子节点
func newLeafNode(p *interiorNode) *leafNode {
	return &leafNode{
		p: p,
	}
}

func (l *leafNode) find(key int) (int, bool) {
	c := func(i int) bool {
		return l.kvs[i].key >= key
	}
	// 查询
	i := sort.Search(l.count, c)
	// 判断是否 key 已经存在
	if i < l.count && l.kvs[i].key == key {
		return i, true
	}

	return i, false
}

// insert
func (l *leafNode) insert(key int, value string) (int, *leafNode, bool) {
	i, ok := l.find(key)

	if ok {
		l.kvs[i].value = value
		return 0, nil, false
	}
	// 判断叶子节点是否已经填满
	if !l.full() {
		copy(l.kvs[i+1:], l.kvs[i:l.count])
		l.kvs[i].key = key
		l.kvs[i].value = value
		l.count++
		return 0, nil, false
	}
	// 获取分裂出新的节点
	next := l.split()
	// 判断插入位置：新的节点或者旧节点
	if key < next.kvs[0].key {
		l.insert(key, value)
	} else {
		next.insert(key, value)
	}
	// 返回分裂节点的第一个key
	return next.kvs[0].key, next, true
}

// 叶子节点分裂过程
func (l *leafNode) split() *leafNode {
	// 申请一个右节点
	next := newLeafNode(nil)
	// 将原始节点的右半部分复制过去
	copy(next.kvs[0:], l.kvs[l.count/2+1:])
	// 初始化原始节点的右半部分
	l.initArray(l.count/2 + 1)
	// 设置右节点的参数
	next.count = MaxKV - l.count/2 - 1
	next.next = l.next
	// 重新设置原始节点的参数
	l.count = l.count/2 + 1
	l.next = next
	// // 设置右节点在 map 中的 key
	// next.setKey(next.kvs[0].key)
	// 返回右节点指针
	return next
}

// 判断是否达到 key 数目上限 MaxKV
func (l *leafNode) full() bool { return l.count == MaxKV }

// 返回父节点，中间节点
func (l *leafNode) parent() *interiorNode { return l.p }

// 设置父中间节点
func (l *leafNode) setParent(p *interiorNode) { l.p = p }

// 获取叶子结点存储的元素数目
func (l *leafNode) countNum() int { return l.count }

// 初始化数组从 num 起的元素为空结构
func (l *leafNode) initArray(num int) {
	for i := num; i < len(l.kvs); i++ {
		l.kvs[i] = kv{}
	}
}
