package skiplist
//https://github.com/ryszard/goskiplist
//包skiplist实现基于跳过列表的地图和集合。
//作者：Ric Szopa（Ryszard）<ryszard.szopa@gmail.com>
//跳过列表是一种数据结构，可以代替
//平衡的树木。跳过列表使用概率平衡而不是
//严格执行平衡，因此算法
//在跳过列表中插入和删除要简单得多，并且
//比平衡树的等效算法快很多。
//https://upload-images.jianshu.io/upload_images/19063731-684743ff91121d5f.jpeg?imageMogr2/auto-orient/strip|imageView2/2/w/1142/format/webp
//跳过列表最初是在William Pugh（1990年6月）中描述的。“跳过
//清单：平衡式的概率替代方案

const SkipList_MaxLevel = 32
const p = 0.25

//存储键值对的容器
type node struct {
	forward    []*node //前面的节点集合
	backward   *node	//后一个节点
	key, value interface{}
}
//next返回skiplist中包含n的下一个节点
func (n *node)next()*node{
	if len(n.forward) == 0 {
		return nil
	}
	return n.forward[0]
}

//SkipList是一种类似于Map的数据结构，可保持有序 键/值对的集合。
// 插入，查找和删除是 所有O（log n）操作。
// SkipList可以有效地存储多达 2 ^ MaxLevel个项目。
type SkipList struct {
	lessThan func(l, r interface{}) bool
	header   *node
	footer   *node
	length   int
	//增加 MaxLevel 是安全的，可以容纳更多的元素
	//如果您降低了 MaxLevel，并且跳过列表已经包含了高级别的节点，
	//那么有效的 MaxLevel 将是新的 MaxLevel 中较大的节点和最高级别的节点。
	//MaxLevel 等于0的 SkipList 等价于标准链表，并且不具备跳过列表的任何好的属性(可能不是您想要的)。
	MaxLevel int
}

//NewCustomMap返回一个新的SkipList，它将使用lessThan作为比较函数。
// lessThan应该在要与SkipList一起使用的键上定义线性顺序。
func NewCustomMap(lessThan func(l, r interface{}) bool) *SkipList {
	return &SkipList{
		lessThan: lessThan,
		header: &node{
			forward: [] *node{nil},
		},
		MaxLevel: SkipList_MaxLevel,
	}
}

//有序接口是可以通过LessThan方法线性排序的接口，由此该实例被认为小于其他实例。 此外，在使用==和！=进行比较时，有序实例应具有正确的行为。
type Ordered interface {
	LessThan(other Ordered) bool
}

//NewCustomMap返回一个新的SkipList，它将使用lessThan作为比较函数。 lessThan应该在要与SkipList一起使用的键上定义线性顺序。
func New() *SkipList {
	comparator := func(left, right interface{}) bool {
		//左边的数据小于右边
		return left.(Ordered).LessThan(right.(Ordered))
	}
	return NewCustomMap(comparator)

}

//getPath
//把全部可能包含key的节点路径，加入到update 数组
// 使用构成可能包含key的节点路径的节点填充到update数组。
//多少个层级就会有多少个元素在update数组中
//候选节点将被返回。 如果update为nil，它将被保留（仍将返回候选节点）。
//如果候选节点太多了，数组太小，将会panic ，联想keys *的使用
func (s *SkipList) getPath(current *node, update []*node, key interface{}) *node {
	depth := len(current.forward) - 1 //深度和forward的关系 可以看图
	//类似二分查找的过程 。找比key小的节点
	for i := depth; i >= 0; i-- {//往深处遍历
		for current.forward[i] != nil && s.lessThan(current.forward[i].key,key) {
			current = current.forward[i] //定位到比key小的节点
		}
		if update != nil {
			update[i] = current
		}
	}
	return current.next() //返回下一个节点
}
//查找key对应的值
func (s *SkipList) Get(key interface{})(value interface{},ok bool){
	node := s.getPath(s.header,nil,key)
	if node == nil || node.key != key {
		return nil ,false
	}
	return node.key,true
}
func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}
//返回有效的最大层级 即限制层级
func (s *SkipList)effectiveMaxLevel() int{
	return maxInt(s.getlevel(),s.MaxLevel)
}
//返回层数
func (s *SkipList) getlevel()int {
	return len(s.header.forward) -1
}
func (s *SkipList) Set(key, value interface{}) {
	if key == nil {
		panic("skiplist: set nil is not support")
	}
	//首次设置level 为0 ，
	updateNodes := make([]*node,s.MaxLevel+1,s.effectiveMaxLevel()+1)
	//找到候选节点;直接更新value
	expectNode := s.getPath(s.header,updateNodes,key)
	if expectNode != nil && expectNode.key == key {
		expectNode.value = value
		return
	}

	newLevel := 0 //这里需要实现随机的层级，避免插入同一层级造成性能下降，
	// 需要在所有层级，均可以添加新节点，可以想象成 三峡移民分配到不同的省市区县镇，
	// 要是都移到北京节点，

	//新增一个层级 节点即为头节点
	if currenLevel := s.getlevel();newLevel > currenLevel {
		//在update中没有更高的level的指针。header 应该在哪里，
		//还可以添加更高level的link 到header
		for i:= currenLevel +1 ;i<= newLevel;i++{
			updateNodes = append(updateNodes,s.header)
			s.header.forward = append(s.header.forward,nil)
		}
	}
	//插入到更浅一层
	newNodeB := &node{
		forward: make([]*node, newLevel+1),
		key:     key,
		value:   value,
	}
	//添加值到链表中，
	//a c之间插入b

	//一、将b前置节点为a
	if previous := updateNodes[0];previous.key != nil  {
		newNodeB.backward = previous
	}
	//二 处理索引 新节点后边追加已有节点 ,这里可以根据图方便理解
	var previousNodeA, previousNodeC *node
	for level:=0;level <= newLevel ;level++{
		//1、找到节点 a
		previousNodeA = updateNodes[level]
		//2 通过节点a 找到节点c
		previousNodeC = previousNodeA.forward[level]
		//3、 b指向（a的下一个节点c）
		newNodeB.forward[level] = previousNodeC
	}
	//长度加1
	s.length++
	//三、给节点C增加前置节点 B，因为只需要加一次，因此不能放到循环中
	if newNodeB.forward[0] != nil {
		if newNodeB.forward[0].backward != newNodeB {
			newNodeB.forward[0].backward = newNodeB
		}
	}
	//如果新加节点比footer大，替换footer
	if s.footer == nil || s.lessThan(s.footer,key){
		s.footer = newNodeB
	}

	s.length++

}
