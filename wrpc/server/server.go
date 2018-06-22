package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"wrpc/constant"
	"wrpc/entity/server"
	"wrpc/werror"
)

type Server struct {
	Conn  net.Conn
	// 服务端应该不需要主动关闭连接
	// 此处先注释不关闭，由客户端进行关闭
	//close chan uint8
}

func StartServerForTCP(host string)  {
	addr, err := net.ResolveTCPAddr("tcp", host)
	werror.CheckError(err)
	lister, err := net.ListenTCP("tcp", addr)
	defer lister.Close()
	for {
		conn, err := lister.Accept()
		s := &Server{
			Conn:  conn,
			//close: make(chan uint8),
		}
		werror.CheckError(err)
		go s.Handle()
	}
}

var MethodPool map[string]*server.Service

// 服务端应该不需要主动关闭连接
// 此处先注释不关闭，由客户端进行关闭
func (s *Server) Close() {
	//select {
	//case <-s.close:
	//	fmt.Println("正常关闭")
	//	close(s.close)
	//	//s.Conn.Close()
	//case <-time.After(constant.TIME_OUT):
	//	fmt.Println("超时关闭")
	//	//s.Conn.Close()
	//	close(s.close)
	//}
}

// 接收到请求后，进行统筹处理
func (s *Server) Handle() {
	defer s.Close()
	// 用于存放接收到信息
	var head *server.CallInfo = new(server.CallInfo)
	// 调用结果通知
	head.Reply = make(chan server.Response)
	// 读取请求
	go s.readRequest(head)
	// 响应请求
	go s.sendResponse(head)
}

// 读取请求（采用json格式进行交互）
func (s *Server) readRequest(head *server.CallInfo) {
	bufferReader := bufio.NewReader(s.Conn)
	for {
		// 按行读取数据。客户端发送消息时，一条调用必须在一行内（协议）
		content, err := bufferReader.ReadString(constant.END_SIGN)
		if len(content) > 0 && err == nil {
			fmt.Println(fmt.Sprintf("通过TCP监听到的数据为：%s", content))
		} else {
			continue
		}
		// 将接收到的信息，转换成对应的结构体
		err = json.Unmarshal([]byte(content), &head)
		// 如果出错，跳出读取循环，并关闭当前连接
		if err != nil {
			reply := server.Response{
				Code:    constant.FAILED,
				Message: constant.FAILED_MSG,
				Data:    nil,
			}
			head.Reply <- reply
			break
		}
		// 调用具体的方法
		go invoke(head)
		break
	}
}

// 响应请求
func (s *Server) sendResponse(h *server.CallInfo) {
	// 响应完之后，关闭连接
	//defer func(){s.close <-0}()
	rtn := <-h.Reply
	close(h.Reply)
	// 如果调用过程出错，则直接返回错误信息，不返回数据
	if rtn.Code == constant.SUCCESS {
		var reply []interface{}
		if v, ok := rtn.Data.([]reflect.Value); ok {
			for i := 0; i < len(v); i++ {
				reply = append(reply, v[i].Interface())
			}
		}
		rtn.Data = reply
	}
	data, _ := json.Marshal(rtn)
	// 写入时，必须以换行为结尾
	data = append(data, []byte("\n")...)
	// 写入数据
	s.Conn.Write(data)
}

// 注册方法
// 将服务名称及其可导出方法存放在一个map中
// 作用：一：可以进行拦截作用，对外只提供开放的方法
// 二：go语言中，使用反射调用方法时，必须知道具体的服务实例，才能调用（不如java强大）
// 注：如果存在多个同服务名且同方法名的方法，则需要起个别名
func Register(rcvr interface{}, name string) {
	// 将指针方法也注册进来
	service := reflect.TypeOf(rcvr)
	if service.Kind() != reflect.Ptr {
		service = reflect.PtrTo(service)
	}

	if len(MethodPool) == 0 {
		MethodPool = make(map[string]*server.Service, service.NumMethod())
	}

	// 将所有可导出方法存放进map中
	for i := 0; i < service.NumMethod(); i++ {
		m := service.Method(i)

		s := &server.Service{
			Method:        m,
			OriginService: reflect.ValueOf(rcvr),
		}

		// key的格式：{服务名}.{方法名}
		MethodPool[fmt.Sprintf("%s.%s", func() string {
			if len(name) == 0 {
				return reflect.Indirect(reflect.ValueOf(rcvr)).Type().Name()
			} else {
				return name
			}
		}(), m.Name)] = s
	}
}

// 调用具体的方法
func invoke(head *server.CallInfo) {
	var code, message string
	// 先检查调用的方法是否存在
	s, ok := MethodPool[fmt.Sprintf("%s.%s", head.ServiceName, head.MethodName)]
	if ok { // 存在，继续操作
		// 组装调用所需的参数
		args, err := getParam(s, head)
		// 参数转换错误，直接返回
		if err != nil {
			code = constant.FAILED
			message = constant.FAILED_MSG
			goto FAILED_RETURN
		}
		// 调用
		v := s.Method.Func.Call(args)

		// 返回
		reply := server.Response{
			Code:    constant.SUCCESS,
			Message: constant.SUCCESS_MSG,
			Data:    v,
		}
		head.Reply <- reply
		return
	}
	code = constant.FAILED
	message = constant.METHOD_NOT_EXIST

FAILED_RETURN:
	// 不存在，直接返回
	reply := server.Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
	head.Reply <- reply
}

// 组装调用方法所需的参数
func getParam(s *server.Service, head *server.CallInfo) ([]reflect.Value, error) {
	var args []reflect.Value = make([]reflect.Value, len(head.Args)+1)
	// 其中第一个参数为：服务信息，具体的struct（reflect调用所必需）
	args[0] = s.OriginService
	// 取出其他的参数，并将其转换成对应方法的实体
	for i := 1; i < len(head.Args)+1; i++ {
		var t reflect.Value
		// 反射创建对应的实例
		if s.Method.Type.In(i).Kind() == reflect.Ptr {
			t = reflect.New(s.Method.Type.In(i).Elem())
		} else {
			t = reflect.New(s.Method.Type.In(i))
		}
		// 将值存放进刚创建的实例中
		err := converToStruct(t, head.Args[i-1])
		// 参数转换错误，直接返回
		if err != nil {
			return nil, err
		}
		args[i] = t
	}
	return args, nil
}

// 将获取的json格式的数据，转换成对应方法的实体
func converToStruct(des reflect.Value, ori interface{}) error {
	i := des.Interface()
	data, err := json.Marshal(ori)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, i)
	if err != nil {
		return err
	}
	return nil
}
