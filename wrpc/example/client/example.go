package main

import (
	"net/http"
	"fmt"
)

type Param struct {
	One, Two int
	SonParam
}

type SonParam struct {
	SOne, STwo int
}

func (e *Example) Add(w http.ResponseWriter, r *http.Request){
	// 调用远程方法
	var rply int
	err := e.C.Call("GateWay", "Add3", []interface{}{Param{One:1, Two:2}}, &rply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("add3的结果：%v", rply))
}