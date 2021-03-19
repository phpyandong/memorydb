package uuid

import (
	"io"
	"fmt"
	"crypto/rand"
)

// GenerateRandomBytes用于生成给定大小的随机字节。

// GenerateRandomBytesWithReader用于生成从给定阅读器读取的给定大小的随机字节。
func GenerateRandomBytesWithReader(size int,reader io.Reader) ([]byte,error){
	if reader == nil {
		return nil ,fmt.Errorf("provided reader is nil")
	}
	buf := make([]byte,size)
	if _,err := io.ReadFull(reader,buf);err != nil{
		return nil, fmt.Errorf("failed to read radom bytes :%v",err)
	}
	return buf,nil
}
const bufLen = 16
// GenerateUUID用于生成随机UUID
func GenerateUUID ()(string, error){
	return GenerateUUIDWithReader(rand.Reader)
}

// GenerateUUIDWithReader用于生成具有给定Reader的随机UUID
func GenerateUUIDWithReader(reader io.Reader)(string,error){
	if reader == nil {
		return "",fmt.Errorf("provided reader is nil")
	}
	buf, err := GenerateRandomBytesWithReader(bufLen,reader)
	if err != nil {
		return "",err
	}
	return FormatUUID(buf)
}
//4b9f7904-3ae4-9acf-a73f-a22eb1f738e0

func FormatUUID (buf []byte)(string,error){
	if buflen := len(buf) ; buflen  != bufLen {
		return "",fmt.Errorf("wrong length byte slice (%d)",buflen)
	}
	fmt.Println(len("x8"))

	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",//十六进制显示
		buf[0:4],//4
		buf[4:6],
		//两个10进制数单字节最大256
		// 249 转为十六进制 f9
		// 221 转为十六进制 dd
		// 一个长度位的十进制转换成16进制 要用2个字节
		buf[6:8],//2
		buf[8:10],//2
		buf[10:16],//6
	) ,nil
}
//func ParseUUID(uuid string)([]byte,error){
//	//1个16进制位 转成二进制位 是2个长度的。
//	// 一个byte = 8 位 二进制 2个16进制位 = 1byte
//	//16个16进制位 = 8byte
//	if len(uuid) !=2 * uuidLen
//
//	}
//}