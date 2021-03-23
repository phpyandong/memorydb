package set

import (
	"sync"
)

//https://github.com/deckarep/golang-set/
//包mapset实现了一个简单而通用的集合集合。
//存储在其中的项目是无序且唯一的。它支持
//典型的集合操作：成员资格测试，交集，联合，
//差异，对称差异和克隆。
//
//包mapset提供了Set的两种实现
// 界面。默认实现对于并发是安全的
//访问，但是还提供了非线程安全的实现
//可以从轻微的速度改进中受益的程序
//可以通过其他方式强制执行互斥。
//包映射集

// Set是mapset包提供的主要接口。它
//表示无序的数据集和大量的
//可以应用于该集合的操作。
type Set interface {
	//将数据添加到集合，返回是否添加成功
	Add(i interface{}) bool
	//set 是否包含i
	Contains(i ...interface{}) bool
}
func NewSet(s ...interface{}) Set{
	set := NewThreadUnsafeSet()
	for _,item := range s{
		set.Add(item)
	}
	return &set
}

type threadUnsafe map[interface{}]struct{}

func NewThreadUnsafeSet() threadUnsafe{
	return make(threadUnsafe)

}
func (set *threadUnsafe) Add(i interface{}) bool {
	_,found := (*set)[i]
	if found {
		return false// 已经存在
	}
	(*set)[i] = struct{}{}
	return true
}

func (set *threadUnsafe) Contains(i ...interface{}) bool {
	for _,val := range i{
		if _,ok := (*set)[val];!ok{
			return false //有一个不存在就返回false
		}
	}
	return true
}

type SafeSet struct {
	s threadUnsafe
	sync.RWMutex
}

func NewUnSafeSet() SafeSet {
	return SafeSet{
		s: NewThreadUnsafeSet(),
	}
}
func (set *SafeSet) Add(i interface{}) bool {
	set.Lock()
	ret := set.s.Add(i)
	set.Unlock()
	return ret
}

func (set *SafeSet) Contains(i ...interface{}) bool {
	set.RLock()
	ret := set.s.Contains(i)
	set.Unlock()
	return ret
}


