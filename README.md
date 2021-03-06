## 基于 map 的 B+ Tree 设计方案

### B 树结构
  http://taop.marchtea.com/03.02.html
### B+ 树的改进
#### B+ 树变更
```
1. 叶子结点存储了所有的关键字信息，并且通过右指针形成链表，可以做到从小到大的顺序遍历；
2. 所有中间节点可以看作索引，结点中仅含有其子树根结点中最大（或最小）关键字。
```
#### B+ 树优势
```
1. B+ 树中间节点只是作为索引来用，占用空间小，读一次磁盘可以加载更多的索引，减少了查询需要的 IO 次数；
2.任何关键字的查找必须走一条从根结点到叶子结点的路，所有关键字查询的路径长度相同，导致每一个数据的查询效率相当；
3.有序数组链表简化了树的遍历操作。
```

### 数据结构定义
#### 建议
	1.所有数据结构应该存储在 kv 对里，整个结构对外表现为一个 map，B+ Tree 是为了减少 map 上 key 的访问次数和封锁粒度而存在的
	2.允许有序扫描
#### 定义
```
	// B+ 树数据结构；限制只能存储正数
	type BTree map[int]node

```
### 核心接口定义
```
	// node 接口设计
	type node interface {
		// 确定元素在节点中的位置
		find(key int) (int, bool)
		// 获取父亲节点
		parent() *interiorNode
		// 设置父亲节点
		setParent(*interiorNode)
		// 是否达到最大数目限制
		full() bool
		// 元素数目统计	
		countNum() int
	}
	// Tree API
	// 创建自由一个父亲节点和叶子节点的 B+ 树
	func NewBTree() *BTree 
	// 返回 B+ Tree 存储的元素数目
	func (bt *BTree) Count() int 
	// 返回根结点
	func (bt *BTree) Root() node 
	// 返回最左侧叶子结点
	func (bt *BTree) First() node 
	// 返回由叶子结点指针构成的数组，从最左侧开始依次追加
	func (bt *BTree) Values() []*leafNode
	// 在 B+ 树中，插入 key-value
	func (bt *BTree) Insert(key int, value string)
	// 搜索： 找到，则返回 value ，否则返回空value
	func (bt *BTree) Search(key int) (string, bool) 

```
### 进阶
```
	1. parent 指针：包含了指向parent的指针，每次分裂要修改大量子节点的父指针（可以考虑 Stack 解决向上递归的问题） 
	2. 支持高并发修改和查询，但是锁的粒度要小，仅限于相关的页面／节点级别（考虑到并发场景下的节点分裂，Stack 里面可能也需要检查右指针）
	3. 排序算法向上暴露
	4. Key 和 Value 数据类型是 Interface
	5. 阅读 Pgsql 的 Readme 文档和两篇关键论文（google pgsql/src/backend/access/nbtree/README）
```

<!--### 与普通实现方案的对比-->
