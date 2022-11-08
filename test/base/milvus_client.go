package base

import (
	"context"
	"fmt"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MilvusClient struct {
	
}


func before() (client.Client, error) {
	fmt.Println("start connect")
	return client.NewGrpcClient(context.Background(), "10.100.31.107:19530")
}

func after(err error, args ...interface{}) {
	if err != nil {
		fmt.Errorf("Error with q%", err)
	}
	//if len(args) == 0 {
	//
	//}
	fmt.Println(args)
	//for _, v := range args {
	//	fmt.Printf("[ApiResponse]: %q", v)
	//}
}

func CreateAlias(collName string, alias string)  error {
	c, _ := before()
	err := c.CreateAlias(context.Background(), collName, alias)
	after(err)
	return err
}

func ListCollections() ([]*entity.Collection, error) {
	c, _ := before()
	res, err := c.ListCollections(context.Background())
	defer after(err, res)
	return res, err
}

func Play1(a int) int {
	return a+1
}
func Play2(a int, b int) (int, int) {
	return a+1, b+1
}
