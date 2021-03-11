package memorydb

//Indexer是用于定义索引的接口。使用索引以便在MemDB表中高效地查找对象。
//实现SingleIndexer或MultiIndexer之一
//索引器主要负责将查找键返回为 []byte。
//字节片是基础数据存储中的关键数据。
type Indexer interface {
	//调用FromArgs可以从参数列表中构建确切的索引键。
	FromArgs(args ...interface{})
}
//SingleIndexer是一个接口，用于定义为每个对象生成单个值的索引
//单索引
type SingleIndexer interface {
	//FromObject从对象中提取索引值。返回值分别是是否找到索引值、索引值和在提取索引值时的任何错误。
	FromObject(raw interface{}) (bool, []byte, error)
}

//MultiIndexer是一个接口，用于定义为每个对象生成多个值的索引。每个值都存储为一个单独的索引，指向同一个对象。
//例如，提取人名的姓和名并允许基于其中任何一个进行查找的索引就是MultiIndexer。本例的FromObject将分割姓和名，并将两者作为值返回。
//多索引
type MultiIndexer interface {
	//FromObject从对象中提取索引值。除了可以有多个索引值之外，返回值与SingleIndexer相同。
	FromObject(raw interface{}) (bool, [][]byte, error)
}

//StringFieldIndex用于使用反射从对象中提取字段，并在该字段上构建索引。
type StringFieldIndex struct {
	Field string
	Lower bool
}
func (str *StringFieldIndex) FromObject(raw interface{}) (bool,[]byte,error){
	panic("im")
}
func (str *StringFieldIndex) FromArgs(args ...interface{}) {
	panic("implement me")
}

