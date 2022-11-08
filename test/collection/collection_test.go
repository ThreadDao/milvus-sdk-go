package collection

import (
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/test/base"
	"testing"
)


func TestOne(t *testing.T) {
	collections,_ := base.ListCollections()
	for _, collection := range collections {
		fmt.Println(collection.Name)
	}
}

func TestTmp(t *testing.T)  {
	y1 := base.Play1(2)
	fmt.Println(y1)

	y2, y3 := base.Play2(3, 4)
	fmt.Println(y2)
	fmt.Println(y3)
}
