package bplusmtree

// B+ 树数据结构；限制只能存储正数
type BTree struct {
	root     *interiorNode // 根节点
	first    *leafNode // 最左侧的叶子节点
	leaf     int // 叶子数目
	interior int // 中间节点数目
	height   int // 树的高度
}

// 创建自由一个父亲节点和叶子节点的 B+ 树
func newBTree() *BTree {
	leaf := newLeafNode(nil)
	r := newInteriorNode(nil, leaf)
	leaf.p = r
	return &BTree{
		root:     r,
		first:    leaf,
		leaf:     1,
		interior: 1,
		height:   2,
	}
}

// 返回 B+ 树的第一个叶子节点
func (bt *BTree) First() *leafNode {
	return bt.first
}

// 在 B+ 树中，插入 key-value
func (bt *BTree) Insert(key int, value string) {
	// 确定插入的位置，是一个叶子节点
	_, oldIndex, leaf := search(bt.root, key)
	// 获取叶子节点的父亲节点，中间节点
	p := leaf.parent()
	// 插入叶子节点，返回是否分裂
	mid, bump := leaf.insert(key, value)
	// 未分裂，则直接返回
	if !bump {
		return
	}
	// 分裂则继续插入中间节点
	var midNode node
	midNode = leaf
	// 设置父亲节点指向分裂出的子（叶子）节点
	p.kcs[oldIndex].child = leaf.next
	// 新分裂出的节点设置该中间节点为父亲节点
	leaf.next.setParent(p)
	// 赋值，获取该中间节点和其父节点
	interior, interiorP := p, p.parent()
	// 迭代向上判断父亲节点是否需要分裂
	for {
		var oldIndex int
		var newNode *interiorNode
		// 判断是否到达根节点
		isRoot := interiorP == nil
		// 未到达根节点，在父亲节点中查询元素的位置
		if !isRoot {
			oldIndex, _ = interiorP.find(key)
		}
		// 将叶子节点分裂后产生的中间元素同时传给父亲中间节点，并传入分裂后的原始叶子节点
		// 同时返回分裂后产生的中间节点和中间元素的 key
		mid, newNode, bump = interior.insert(mid, midNode)
		// 未分裂，直接返回
		if !bump {
			return
		}

		if !isRoot {
			// 未到达根节点，将元素插入父亲节点，基本过程同上
			interiorP.kcs[oldIndex].child = newNode
			newNode.setParent(interiorP)

			midNode = interior
		} else {
			// 到达根节点，根节点上移，并插入原始中间节点
			bt.root = newInteriorNode(nil, newNode)
			newNode.setParent(bt.root)

			bt.root.insert(mid, interior)
			return
		}
		// 赋值，获取该中间节点的父亲节点和其父亲的父节点
		interior, interiorP = interiorP, interior.parent()
	}
}

// 搜索： 找到，则返回 value ，否则返回空value
func (bt *BTree) Search(key int) (string, bool) {
	kv, _, _ := search(bt.root, key)
	if kv == nil {
		return "", false
	}
	return kv.value, true
}

// 具体搜索过程
func search(n node, key int) (*kv, int, *leafNode) {
	curr := n
	oldIndex := -1

	for {
		switch t := curr.(type) {
		// 叶子节点，返回命中节点或者可插入位置
		case *leafNode:
			i, ok := t.find(key)
			if !ok {
				return nil, oldIndex, t
			}
			return &t.kvs[i], oldIndex, t
		// 中间节点迭代查询
		case *interiorNode:
			i, _ := t.find(key)
			curr = t.kcs[i].child
			oldIndex = i
		default:
			panic("")
		}
	}
}
