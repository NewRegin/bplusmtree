package bplustree

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
	kvs   kvs // 存储元素
	count int // 实际存储元素数目
	next  *leafNode // 右边第一个叶子节点（右指针）
	p     *interiorNode // 父亲节点，中间节点
}

// 创建新的叶子节点
func newLeafNode(p *interiorNode) *leafNode {
	return &leafNode{
		p: p,
	}
}

// find finds the index of a key in the leaf node.
// If the key exists in the node, it returns the index and true.
// If the key does not exist in the node, it returns index to
// insert the key (the index of the smallest key in the node that larger
// than the given key) and false.
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
func (l *leafNode) insert(key int, value string) (int, bool) {
	i, ok := l.find(key)

	if ok {
		//log.Println("insert.replace", i)
		l.kvs[i].value = value
		return 0, false
	}
	// 判断叶子节点是否已经填满
	if !l.full() {
		copy(l.kvs[i+1:], l.kvs[i:l.count])
		l.kvs[i].key = key
		l.kvs[i].value = value
		l.count++
		return 0, false
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
	return next.kvs[0].key, true
}

// 叶子节点分裂过程
func (l *leafNode) split() *leafNode {
	// 申请一个右节点
	next := newLeafNode(nil)
	// 将原始节点的右半部分复制过去
	copy(next.kvs[0:], l.kvs[l.count/2+1:])
	// 设置右节点的参数
	next.count = MaxKV - l.count/2 - 1
	next.next = l.next
	// 重新设置原始节点的参数
	l.count = l.count/2 + 1
	l.next = next
	// 返回右节点指针
	return next
}
// 判断是否达到 key 数目上限 MaxKV
func (l *leafNode) full() bool { return l.count == MaxKV }
// 返回父节点，中间节点
func (l *leafNode) parent() *interiorNode { return l.p }
// 设置父中间节点
func (l *leafNode) setParent(p *interiorNode) { l.p = p }
