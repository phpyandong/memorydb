package safemap

import (
	"testing"
	"fmt"
	"time"
)

func TestNewConcurrencyMap(t *testing.T) {
	cmap := NewConcurrencyMap()
	cmap.Set("seiya1","seiya1res")
	cmap.Set("seiya2","seiya2res")
	//cmap.Set("seiya3","seiya3")
	//cmap.Set("seiya4","seiya4")
	fmt.Println(cmap.Get("seiya1"))
	fmt.Println(cmap.Get("seiya2"))
	//fmt.Println(cmap.Get("seiya3"))
	//fmt.Println(cmap.Get("seiya4"))
	//fmt.Println(cmap.Remove("seiya4"))
	ch := cmap.Elements()
	select {
	case mes:=<-ch:

		fmt.Println("eles",mes)
	}

	time.Sleep(10*time.Second)

}