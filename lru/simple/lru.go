package simple

import (
	"container/list"
)

type Key interface{}
type entry struct {
	key   Key
	value interface{}
}

//缓存是LRU缓存。并发访问是不安全的。
type Cache struct {
	//MaxEntries是逐出项之前缓存项的最大数目。零意味着没有限制。
	MaxEntries int
	//OnEvicted（可选）指定从缓存中清除条目时要执行的回调函数。
	OnEvicted func(key Key, value interface{})
	lists     *list.List
	cacheMap     map[interface{}]*list.Element //使用map方式快速查找到list中元素的指针
}

//如果maxEntries为零，则缓存没有限制，并且假定逐出是由调用方完成的。
func NewCache(maxEntries int) *Cache {
	return &Cache{
		MaxEntries: maxEntries,
		lists:      list.New(),
		cacheMap:      make(map[interface{}]*list.Element),
	}
}
func (c *Cache) Get(key Key)(value interface{},ok bool){
	if c.cacheMap == nil {
		return
	}
	if ele ,hit := c.cacheMap[key]; hit{
		c.lists.MoveToFront(ele)
		return ele.Value.(*entry).value,true
	}
	return
}
func (c *Cache) Add(key Key, value interface{}) {
	if c.cacheMap == nil {
		c.cacheMap = make(map[interface{}]*list.Element)
		c.lists = list.New()
	}
	//如果有值
	if listElement, ok := c.cacheMap[key]; ok {
		c.lists.MoveToFront(listElement)
		listElement.Value.(*entry).value = value
		return
	}
	//如果没有值;
	listNewItem := c.lists.PushFront(&entry{key: key, value: value})
	c.cacheMap[key] = listNewItem
	//每次添加时，判读是否达到配额， 触发删除旧的缓存
	if c.MaxEntries != 0 && c.lists.Len() > c.MaxEntries {
		c.RemoveOldest()
	}

}
//删除list中的尾部元素
func (c *Cache) RemoveOldest() {
	if c.cacheMap ==nil {
		return
	}
	ele := c.lists.Back()//找到最后一个元素
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *list.Element){
	c.lists.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cacheMap,kv.key) //删除map 并发不安全
	if c.OnEvicted != nil {
		//移除后的回调函数.可以实现监控某个键
		c.OnEvicted(kv.key,kv.value)
	}
}

func (c *Cache) Len() int{
	if c.cacheMap == nil {
		return 0
	}
	return c.lists.Len()
}

func (c *Cache) Clear(){
	if c.OnEvicted != nil {
		for _, e := range c.cacheMap{
			kv := e.Value.(*entry)
			c.OnEvicted(kv.key,kv.value)
		}
	}
	c.lists = nil
	c.cacheMap = nil
}