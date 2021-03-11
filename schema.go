package memorydb
////表是数据库中的表的集合。键是表名，必须与表表模式中的名称匹配。
type DBSchema struct {
	Tables map[string]*TableSchema
}
//TableSchema是表的chema。
type TableSchema struct {
	//表名。它必须与DBSchema中的表映射中的键相匹配。
	Name	string
	//索引是用于查询该表的索引集。键是索引的唯一名称，必须与IndexSchema中的名称匹配。
	Indexes map[string]*IndexSchema
}

//IndexSchema是索引的schema。索引定义了如何查询一个表。
type IndexSchema struct {
	//索引的名称。
	//这在一个索引集的表中必须是唯一的。
	//这必须与表模式的索引映射中的键相匹配
	Name string
	AllowMissing bool
	Unique 		bool //是否唯一
	Indexer 	Indexer

}