package main

import (
	"fmt"
)

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
	fmt.Println("接收到参数，并返回", p, p.One+p.Two)
	return p.One + p.Two
}

func (gw GateWay) Add4() int {
	return 1 + 2
}
