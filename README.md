# memorydb
本地实时内存数据库，支持事务，及lru

提供' memorydb '包，实现了一个基于不可变基数树的简单内存数据库。数据库提供了原子性、一致性和与ACID的隔离性。
Atomicity, Consistency, and Isolation from ACID.
因为它在内存中，所以不能提供持久性。
数据库用一个模式来实例化，该模式指定存在的表和索引，并允许执行事务。

数据库提供以下信息:

*多版本并发控制(MVCC) -
通过利用不可变的基数树，数据库能够在不锁定的情况下支持任意数量的并发读取，并允许写入操作。

*事务支持——数据库支持丰富的事务，可以插入、更新或删除多个对象。
事务可以跨多个表，并自动应用。
数据库在ACID术语中提供了原子性和隔离性，因此在提交之前更新是不可见的。

*丰富的索引——表可以支持任意数量的索引，可以是简单的单个字段索引，也可以是更高级的复合字段索引。
某些类型(如UUID)可以有效地从字符串压缩到字节索引，以减少存储需求。

*丰富的索引——表可以支持任意数量的索引，可以是简单的单个字段索引，也可以是更高级的复合字段索引。
某些类型(如UUID)可以有效地从字符串压缩到字节索引，以减少存储需求。


*Watches-调用者可以填充一个watch set 作为查询的一部分，这可以用来检测当修改数据库影响查询结果。
这使得调用者可以很容易地在非常普遍的情况下观察数据库中的更改


Example
=======

Below is a [simple example](https://play.golang.org/p/gCGE9FA4og1) of usage

```go
// Create a sample struct
type Person struct {
	Email string
	Name  string
	Age   int
}

// Create the DB schema
schema := &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		"person": &memdb.TableSchema{
			Name: "person",
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "Email"},
				},
				"age": &memdb.IndexSchema{
					Name:    "age",
					Unique:  false,
					Indexer: &memdb.IntFieldIndex{Field: "Age"},
				},
			},
		},
	},
}

// Create a new data base
db, err := memdb.NewMemDB(schema)
if err != nil {
	panic(err)
}

// Create a write transaction
txn := db.Txn(true)

// Insert some people
people := []*Person{
	&Person{"joe@aol.com", "Joe", 30},
	&Person{"lucy@aol.com", "Lucy", 35},
	&Person{"tariq@aol.com", "Tariq", 21},
	&Person{"dorothy@aol.com", "Dorothy", 53},
}
for _, p := range people {
	if err := txn.Insert("person", p); err != nil {
		panic(err)
	}
}

// Commit the transaction
txn.Commit()

// Create read-only transaction
txn = db.Txn(false)
defer txn.Abort()

// Lookup by email
raw, err := txn.First("person", "id", "joe@aol.com")
if err != nil {
	panic(err)
}

// Say hi!
fmt.Printf("Hello %s!\n", raw.(*Person).Name)

// List all the people
it, err := txn.Get("person", "id")
if err != nil {
	panic(err)
}

fmt.Println("All the people:")
for obj := it.Next(); obj != nil; obj = it.Next() {
	p := obj.(*Person)
	fmt.Printf("  %s\n", p.Name)
}

// Range scan over people with ages between 25 and 35 inclusive
it, err = txn.LowerBound("person", "age", 25)
if err != nil {
	panic(err)
}

fmt.Println("People aged 25 - 35:")
for obj := it.Next(); obj != nil; obj = it.Next() {
	p := obj.(*Person)
	if p.Age > 35 {
		break
	}
	fmt.Printf("  %s is aged %d\n", p.Name, p.Age)
}
// Output:
// Hello Joe!
// All the people:
//   Dorothy
//   Joe
//   Lucy
//   Tariq
// People aged 25 - 35:
//   Joe is aged 30
//   Lucy is aged 35
```

