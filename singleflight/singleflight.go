package singlefilght

import "sync"


type Group struct {
	mutex sync.Mutex
	mapKeyCall map[string]*Call// 对于每一个需要获取的key有一个对应的call
}
// call代表需要被执行的函数
type Call struct {
	waitgrp sync.WaitGroup //用于阻塞这个调用call的其他请求
	result interface{} 	//函数执行后的结果
	err error //函数执行后的error
}
func(group *Group) Exec(key string,fn func()(interface{},error))(interface{},error	){
	group.mutex.Lock() //给整个group上大锁
	if group.mapKeyCall == nil {
		group.mapKeyCall = make(map[string]*Call)
	}
	//如果获取当前key的函数正在被执行，则阻塞等待执行中的，等待其执行完毕后获取他的执行给结果
	if call ,ok := group.mapKeyCall[key];ok {
		group.mutex.Unlock()
		call.waitgrp.Wait()
		return call.result,call.err
	}
	//如果没有执行过,初始化一个call ,往map中写后解锁
	call := new(Call)
	call.waitgrp.Add(1)//正在执行，进行阻塞读
	group.mapKeyCall[key] = call
	group.mutex.Unlock()

	//执行获取key的函数，并将该结果过赋值给Call
	call.result,call.err = fn()
	call.waitgrp.Done() //执行完毕，不再阻塞读取

	//此时缓存应读本地。可以在group中删除这个key
	group.mutex.Lock()
	delete(group.mapKeyCall,key)
	group.mutex.Unlock()

	return call.result,call.err
}