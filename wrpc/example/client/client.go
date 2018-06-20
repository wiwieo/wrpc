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

func main() {
	conn, err := net.Dial("tcp", "localhost:3333")
	werror.CheckErr(err, 1)
	defer func(){
		fmt.Println("中止")
		conn.Close()
	}()
	c := client.NewClient(conn)
	var rply int
	err = c.Call("GateWay", "Add3", []interface{}{Param{One:1, Two:2}}, &rply)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("add3的结果：", rply)
}
