package snowflake

import (
	"sync"
	"time"
	"fmt"
)
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

	time int64 //41位 毫秒时间戳
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
	n.stepMask = -1 ^ (-1 << StepBits)
	n.timeShift = NodeBits + StepBits //10+12
	n.nodeShift = StepBits
	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	//这里可能是 解决时钟回调的问题？
	n.epoch = curTime.Add( //当前时间加上时间d
		time.Unix(Epoch/1000,(Epoch%1000) * 1000 * 1000 ).Sub(curTime), //
		//根据预设的时间创建一个时间戳，减去当前时间 得到一个差值；如果
		)

	return &n ,nil
}

func (node *Node)Generate(){
	node.mutex.Lock()
	//服务启动的写死了
	now := time.Since(node.epoch).Nanoseconds() /(1000 * 1000)
	if now == node.time {
		fmt.Println("83 等于")
		node.step = (node.step +1) & node.stepMask
		fmt.Println("等于:",node.step)
	}else{
		fmt.Println("83 不等于")

		node.step = 0
	}

	node.time = now
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
}