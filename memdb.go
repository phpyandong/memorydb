package memorydb

import (
	"unsafe"
	"sync"
	"go-immutable-radix"
	"sync/atomic"
	"fmt"
)

//插入到MemDB中的对象不会被复制。
//这是**极其重要的**对象在插入后不能被就地修改，因为它们直接存储在MemDB中。
//即使已经从MemDB中删除了插入的对象，修改它们仍然是不安全的，因为从其他goroutines读取DB的快照可能仍然是旧的。

type MemDB struct {
	schema	*DBSchema
	root 	unsafe.Pointer // *iradix.Tree underneath
	primary bool
	writer  sync.Mutex
}

func NewMemDB(schema *DBSchema) (*MemDB,error){
	//校验
	db := &MemDB{
		schema:schema,
		root:	unsafe.Pointer(iradix.New()),
		primary: true,
	}
	if err := db.initialize(); err != nil{
		return nil,err
	}
}
//initialize用于设置创建后使用的数据库。在分配一个MemDB之后，这个函数只能调用一次。
func (db *MemDB) initialize() error{
	root := db.getRoot()
	var oldRoot interface{}
	var ok bool
	for tableName,tableSchema := range db.schema.Tables {
		for indexName := range tableSchema.Indexes {
			index := iradix.New()
			path := indexPath(tableName,indexName)
			//*Tree, interface{}, bool :=
			root, oldRoot,ok = root.Insert(path,index)
			if !ok {
				panic(fmt.Sprintf("init err %v",oldRoot))
			}
		}
	}
	db.root = unsafe.Pointer(root)
	return nil
}
//ls表名.索引
func indexPath(table,index string) []byte {
	return []byte(table + "." + index)
}
func (db *MemDB) getRoot() *iradix.Tree{
	root := (*iradix.Tree)(atomic.LoadPointer(&db.root))
	return root
}