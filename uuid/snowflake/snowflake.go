package snowflake

import (
	"sync"
	"time"
	"fmt"
)
//| 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
//使用默认设置，这允许每个节点ID每毫秒生成4096个唯一ID。
//https://github.com/bwmarrin/snowflake
//一个非常简单的Twitter雪花生成器。
//解析现有雪花ID的方法。
//将雪花ID转换为其他几种数据类型并返回的方法。
//JSON Marshal / Unmarshal函数可轻松在JSON API中使用雪花ID。
//单调时钟计算可防止时钟漂移。
var (
	//将纪元设置为2010年11月4日UTC的twitter雪花纪元（以毫秒为单位）
	Epoch int64 = 1288834974657
	//NodeBits保存用于Node的位数

	//请记住，Node / Step之间共有22位共享
	NodeBits uint8 = 10
	//StepBits保存用于Step的位数

	//记住，节点/步骤之间总共有22位可共享
	StepBits uint8 = 12
	mu		sync.Mutex
	//nodeMax int64= -1 ^(-1 << NodeBits)  //快速求出最大数
	//nodeMask 	= nodeMax << StepBits//倍数

	//stepMask int64 = -1 ^ (-11 << StepBits)
	//timeShift = NodeBits + StepBits //10+12
	//nodeShift = StepBits //12

)

func main()  {
	
}
//推特的id构成（从最高位往最低位方向）：
//
//1位 ，不用。固定是0
//41位 ，毫秒时间戳
//5位 ，数据中心ID (用于对数据中心进行编码)
//5位 ，WORKERID (用于对工作进程进行编码)
//
type Node struct{
	mutex sync.Mutex
	epoch time.Time

	lasttime int64 //41位 毫秒时间戳
	node int64 // 5位 ，数据中心ID (用于对数据中心进行编码)
	step int64 // 12位 ，序列号。用于同一毫秒产生ID的序列 （自增id）

	nodeMax int64
	nodeMask 	int64
	timeShift uint8
	nodeShift uint8
	stepMask int64
}
func NewNode(node int64)(*Node,error){
	//mu.Lock()
	//nodeMax = -1 ^ (-1 << NodeBits) //10
	//nodeMask = nodeMax << StepBits //12
	//stepMask = -1 ^(-1 << StepBits) //12
	//timeShift = NodeBits + StepBits //10+12
	//nodeShift = StepBits //12
	//mu.Unlock()
	n := Node{}
	n.node = node //节点数量
	n.nodeMax = -1 ^ (-1 << NodeBits)
	n.nodeMask = n.nodeMax << StepBits //12
	n.stepMask = -1 ^ (-1 << StepBits) //12位的最大值
	n.timeShift = NodeBits + StepBits //10+12
	n.nodeShift = StepBits
	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	//这里解决时钟回调的问题
	//https://www.yuque.com/chao77977/golang/nafo0g
	//https://zhuanlan.zhihu.com/p/47754783
	//Wall Clocks，顾名思义，表示墙上挂的钟，在这里表示我们平时理解的时间，
	// 存储的形式是自 1970 年 1 月 1 日 0 时 0 分 0 秒以来的时间戳，
	// 当系统和授时服务器进行校准时间时间操作时，
	// 有可能造成这一秒是2018-1-1 00:00:00，而下一秒变成了2017-12-31 23:59:59的情况。
	// Monotonic Clocks，意思是单调时间的，所谓单调，就是只会不停的往前增长，
	// 不受校时操作的影响，这个时间是自进程启动以来的秒数。
	//如果每隔一秒生成一个Time并打印出来，就会看到如下输出
	//2018-10-26 14:15:50.306558969 +0800 CST m=+0.000401093
	//可以看到m=+后面所显示的数字，就是文档中所说的Monotonic Clocks。
	fmt.Println("时间间隔：",	time.Unix(Epoch/1000,(Epoch%1000) * 1000 * 1000 ).Sub(curTime))
	n.epoch = curTime.Add( //当前时间加上时间d ,差值是根据m=+0.000401093 来计算的保证了时钟的单调递增
		time.Unix(Epoch/1000,(Epoch%1000) * 1000 * 1000 ).Sub(curTime), //
		//根据预设的时间创建一个时间戳，减去当前时间 得到一个差值；如果
		)
	//fmt.Println("epoch:",n.epoch)
	//epoch: 2010-11-04 09:42:54.657 +0800 CST m=-327374333.383055884  服务开始启动后，
	// 利用真实时间的原子钟和epoch的原子钟的差作为当前的时间戳，由于真实时间的单调钟，是单调递增的，
	// 因此差值就是单调递增的
	return &n ,nil
}

func (node *Node)Generate() int64{
	node.mutex.Lock()
	//服务启动的写死了
	fmt.Println("epoch:",node.epoch)
	//单调钟的核心应用
	now := time.Since(node.epoch).Nanoseconds() /(1000 * 1000)
	//老版本
	//	now := time.Now().UnixNano() / 1000000 直接用的系统时间

	fmt.Println("now:",now)
	if now == node.lasttime {
		fmt.Println("83 等于")
		fmt.Printf("step1: %b",node.step+1)
		fmt.Printf("step2: %b",node.stepMask)

		node.step = (node.step +1) & node.stepMask
		fmt.Printf("step3: %b",node.step )
		fmt.Printf("step4:%b",(node.stepMask+1) & node.stepMask) //可以防止超过最大数 )
		//如果超过最大数，结果则为0
		if node.step == 0 {
			//如果当前时间 小于上次的时间 进入死循环,
			//直到过了1毫秒不再循环，相当于sleep 一毫秒，解决序列号用完的问题，
			// 不用sleep的原因可能是避免任何使用普通系统时间的情况
			for now <= node.lasttime {
				now = time.Since(node.epoch).Nanoseconds() /(1000 * 1000)
			}

		}

		fmt.Println("等于:",node.step)
	}else{
		fmt.Println("83 不等于")

		node.step = 0
	}

	node.lasttime = now
	//r :=
	r1 := (now) << node.timeShift //44位的时间戳左移22位  xxxxxxxxx 0000000
	//fmt.Printf("r1:%b\n",r1)
	r2 := (node.node << node.nodeShift) //节点数  左移12位
	//fmt.Printf("r2:%b\n",r2)
	r3 := node.step //现在序列号
	//fmt.Printf("r3:%b\n",r3)

	//fmt.Printf("结果：%b\n",r1|r2|r3)
	fmt.Println("结果：",r1|r2|r3)

	node.mutex.Unlock()
	return r1|r2|r3
}