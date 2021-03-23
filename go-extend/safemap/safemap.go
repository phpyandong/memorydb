package safemap

import (
	"sync"
	"github.com/spaolacci/murmur3"
	"bytes"
	"encoding/binary"
	"hash/fnv"
	"hash"
)
type IConcurrencyMap interface {
	// Get 获取给定键值对应的元素值。若没有对应的元素值则返回nil
	Get(key interface{}) (interface{},error)
	// Set 给指定的key 设置value ,若该键值已存在，则替换
	Set(key interface{},elem interface{}) error
	//Remove 删除给定键值对，并返回旧的元素值，若没有旧的元素，则返回nil
	Remove(key interface{})(bool, error)
	//Elements 获取并发Map中的全部元素
	Elements() <-chan ConcurrencyMap
}
// ConcurrencyElement 存储的元素项
type concurrencyElement struct {
	Key interface{}
	Value interface{}
}
type bucket struct{
	sync.RWMutex
	items map[interface{}]interface{}
}
//ConcurrencyMap 是由多个小的map构成的
type ConcurrencyMap struct {
	sync.RWMutex
	size int64
	pools []*bucket
}
//根据key 进行hash 然后取余，得到桶bucket的编号，将数据存到该bucket上
func (cm *ConcurrencyMap) getBucket(key interface{}) (*bucket,error){
	var v interface{}
	switch key.(type) {
	case string:
		v = []byte(key.(string))
	case int:
		v = int32(key.(int))
	default:
		v =key
	}
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer,binary.LittleEndian,v)
	if err != nil {
		return nil,err
	}
	var hasher hash.Hash32

	//OneAtATime - 354.163715 mb/sec
	//FNV - 443.668038 mb/sec
	//SuperFastHash - 985.335173 mb/sec
	//lookup3 - 988.080652 mb/sec
	//MurmurHash 1.0 - 1363.293480 mb/sec
	//MurmurHash 2.0 - 2056.885653 mb/sec

	if true {
		hasher = fnv.New32()
	}else{
		//更快的hash
		hasher = murmur3.New32()
	}
	hasher.Write(buffer.Bytes())

	defer buffer.Reset()

	return cm.pools[uint(hasher.Sum32()) & uint(cm.size)] ,nil
}
func (cmap *ConcurrencyMap) Get(key interface{}) (interface{}, error) {
	bucket,err := cmap.getBucket(key)
	if err != nil {
		return nil,err
	}
	bucket.RLock()
	v := bucket.items[key]
	bucket.Unlock()
	return v,nil
}

func (cmap *ConcurrencyMap) Set(key interface{}, elem interface{}) error {
	bucket ,err := cmap.getBucket(key)
	if err != nil {
		return err
	}
}

func (ConcurrencyMap) Remove(key interface{}) (bool, error) {
	panic("implement me")
}

func (ConcurrencyMap) Elements() <-chan ConcurrencyMap {
	panic("implement me")
}

const DefaultPoolSize = 1 << 5
//创建并发的Map接口 使用多参数的目的，构造非必传参数的效果，骚操作
//poolSize 分配共享池的大小，默认为32 必须2的整数次幂
func NewConcurrencyMap(poolSizes ...uint) IConcurrencyMap{
	var size uint
	if len(poolSizes) >0 {
		size = poolSizes[0]
	}else{
		size = DefaultPoolSize
	}
	pools := make([]*bucket,size)
	for i:= 0;i< int(size) ;i++{
		pools[i] = &bucket{
			items : make(map[interface{}]interface{}),
		}
	}
	return &ConcurrencyMap{
		size :int64(size),
		pools:pools,
	}
}