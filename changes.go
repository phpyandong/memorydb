package memorydb
//在一个事务中对memDB表执行的一组变化。
type Changes []Change
//Change描述表中对象的变化。
type Change struct {
	Table string
	Before interface{}
	After 	interface{}
	//primaryKey存储来自主索引的原始键值，以便我们可以在同一事务中对同一对象进行多次更新，但不向使用者公开这个实现细节。
	primaryKey []byte
}