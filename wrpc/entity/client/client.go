package client

type Header struct {
	ServiceName string
	MethodName string
	Args []interface{}
	Reply []interface{}
}

type Response struct {
	Data    interface{}
	Code    string
	Message string
}