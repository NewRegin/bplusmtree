package bplusmtree

const (
	// 叶子节点最大元素存储数目
	MaxKV = 255
	// 中间节点最大元素存储数目
	MaxKC = 511
)

// 接口设计
type node interface {
	// 确定元素在节点中的位置
	find(key int) (int, bool)
	// 获取父亲节点
	parent() *interiorNode
	// 设置父亲节点
	setParent(*interiorNode)
	// 是否达到最大数目限制
	full() bool
}
