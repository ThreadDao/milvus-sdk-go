package base

import (
	"fmt"
)

type MilvusClientFunc func(args []interface{}) []interface{}

func wrapper(f MilvusClientFunc) MilvusClientFunc {
	return func(args []interface{}) []interface{} {
		fmt.Println("start")
		c := f(args)  // 执行乘法运算函数
		fmt.Println("end")
		return c  // 返回计算结果
	}
}
