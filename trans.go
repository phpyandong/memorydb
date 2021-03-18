package memorydb

import "go-immutable-radix"
type tableIndex struct {
	Table string
	Index string
}
type Txn struct {
	db *MemDB
	write bool
	rootTxn *iradix.Txn
	after 	[]func()
	// changes用于跟踪事务期间执行的更改。如果在事务开始时为nil，则不会跟踪更改。
	changes 	Changes
	modiried map[tableIndex]*iradix.Txn
}
