package main

import (
	"net"
	"reflect"
	"wrpc/werror"
	"wrpc/server"
)


type header struct {
	ServiceName string
	MethodName  string
	Args        []interface{}
	Reply       chan []reflect.Value
}

func main() {
	server.Register(&GateWay{}, "")
	addr, err := net.ResolveTCPAddr("tcp", "localhost:3333")
	werror.CheckErr(err, 1)
	lister, err := net.ListenTCP("tcp", addr)
	werror.CheckErr(err, 2)
	defer lister.Close()
	for {
		conn, err := lister.Accept()
		werror.CheckErr(err, 3)
		go server.Handle(conn)
	}
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
	return p.One + p.Two
}

func (gw GateWay) Add4() int {
	return 1 + 2
}
