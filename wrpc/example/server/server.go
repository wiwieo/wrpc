package main

import (
	"reflect"
	"wrpc/server"
	_ "net/http/pprof"
	"time"
)


type header struct {
	ServiceName string
	MethodName  string
	Args        []interface{}
	Reply       chan []reflect.Value
}

func main() {
	server.Register(&GateWay{}, "")
	server.StartServerForTCP(":3333")
	//http.ListenAndServe(":8080", nil)
}

type GateWay struct {
}

func (gw GateWay) Add1(one, two *int) int {
	return *one + *two
}

type Param struct {
	One, Two int
	SonParam
}

type SonParam struct {
	SOne, STwo int
}

func (gw GateWay) Add2(p *Param) int {
	return p.SOne + p.STwo
}

func (gw GateWay) Add3(p *Param) int {
	time.Sleep(3*time.Second)
	return p.One + p.Two
}

func (gw GateWay) Add4() int {
	return 1 + 2
}