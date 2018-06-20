package server

import (
	"reflect"
)

type Service struct {
	Method        reflect.Method
	OriginService reflect.Value
}

type CallInfo struct {
	ServiceName string
	MethodName  string
	Args        []interface{}
	Reply       chan Response
}

type Response struct {
	Data    interface{}
	Code    string
	Message string
}