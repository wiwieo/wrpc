package main

import (
	"net"
	"fmt"
	"wrpc/werror"
	"wrpc/client"
)

type Param struct {
	One, Two int
	SonParam
}

type SonParam struct {
	SOne, STwo int
}
//
//func main() {
//	var wg sync.WaitGroup
//	fmt.Println("开始时间：%v", time.Now())
//	for i:=0;i<1;i++{
//		wg.Add(1)
//		go func(wg *sync.WaitGroup, i int){
//			defer wg.Done()
//			conn, err := net.Dial("tcp", "localhost:3333")
//			werror.CheckErr(err)
//			defer func() {
//				fmt.Println("中止")
//				conn.Close()
//			}()
//			c := client.NewClient(conn)
//			var rply int
//			err = c.Call("GateWay", "Add3", []interface{}{Param{One:1, Two:i}}, &rply)
//			if err != nil {
//				fmt.Println(err)
//			}
//			fmt.Println(fmt.Sprintf("第%d个add3的结果：%v", i, rply))
//		}(&wg, i)
//	}
//	wg.Wait()
//	fmt.Println("结束时间：%v", time.Now())
//}

func main() {
	conn, err := net.Dial("tcp", ":3333")
	werror.CheckErr(err)
	defer func() {
		conn.Close()
	}()
	c := client.NewClient(conn)
	var rply int
	err = c.Call("GateWay", "Add3", []interface{}{Param{One:1, Two:2}}, &rply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("add3的结果：%v", rply))
}


