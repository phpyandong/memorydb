package snowflake

import (
	"testing"
)

func TestNodeGenerate(t *testing.T) {
	node,_ := NewNode(5)
	node.Generate()
	node.Generate()
	node.Generate()
	node.Generate()
	node.Generate()
	node.Generate()

	node.Generate()
	node.Generate()
	node.Generate()
	node.Generate()
	node.Generate()
	node.Generate()
	node.Generate()
	node.Generate()
	node.Generate()
	node.Generate()
	node.Generate()
	node.Generate()
/*	fmt.Printf("%b\n",time.Now().Unix())

	fmt.Printf("%b\n",time.Now().Unix() << 22)*/
	//fmt.Printf("%b\n",1<<10)
	//
	////^作二元运算符就是异或，包括符号位在内，相同为0，不相同为1
	//
	////二、 &^运算符
	////
	////作用：将运算符左边数据相异的位保留，相同位清零
	//
	//
	//fmt.Printf("%b\n",(1<< 2))
	//
	//fmt.Printf("%b\n",-1^-1<<2)
	////            -1
	////1100 0000 0000
	//// 011 1111 1111
	//fmt.Printf("%b\n",(1<< 10))
	//fmt.Printf("%b\n",1^ (1<< 10))

}
