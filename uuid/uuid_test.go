package uuid

import (
	"testing"
	"fmt"
	uuid2 "github.com/google/uuid"
)

func TestFormatUUID(t *testing.T) {
	uuid2.New()
	fmt.Printf("%x\n",10)
	fmt.Printf("%b\n",10)
	uuid,_ := GenerateUUID()
	fmt.Println(uuid)

}