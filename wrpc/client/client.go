package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"wrpc/constant"
	"wrpc/entity/client"
	"wrpc/etcdv3"
	"wrpc/etcdv3/balancer"
)

var CONNECTED_POOL sync.Map

type Client struct {
	addr string
	// 不自动关闭连接，由调用使用方主动关闭连接
	//close chan uint8
	e    *etcdv3.Etcdv3
	conn chan net.Conn
	b    balancer.Balancer
	mux  sync.RWMutex
}

func Dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}

func DialTCPConnect(is_used_etcd bool, host string, balanceType int, serviceName string, endpoints string) (*Client, error) {
	if !is_used_etcd && len(host) == 0 {
		err := fmt.Errorf("不使用etcd时，必须指定连接的服务器")
		return nil, err
	}
	var e *etcdv3.Etcdv3
	if is_used_etcd {
		e, _ = etcdv3.NewEtcdv3(&etcdv3.ServerInfo{
			Name:      serviceName,
			Endpoints: strings.Split(endpoints, ","),
			TTL:       10,
			Ctx:       context.Background(),
		})
	}
	return &Client{addr: host,
		e:    e,
		conn: make(chan net.Conn),
		b:    balancer.GetBalance(balanceType),
	}, nil
}

func loadConn(host string) (net.Conn, error) {
	v, ok := CONNECTED_POOL.Load(host)
	if ok {
		c, ok := v.(net.Conn)
		if !ok {
			return nil, fmt.Errorf("系统错误。")
		}
		return c, nil
	} else {
		return nil, fmt.Errorf("不存在可用的连接。")
	}
}

func storeConn(host string) net.Conn {
	conn, err := loadConn(host)
	if err == nil{
		return conn
	}
	v, ok := CONNECTED_POOL.Load(host)
	if !ok {
		conn, err := Dial("tcp", host)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		CONNECTED_POOL.Store(host, conn)
		return conn
	}
	c, ok := v.(net.Conn)
	if ok {
		return c
	}
	return nil
}

func delConn(host string) {
	CONNECTED_POOL.Delete(host)
}

//func NewClient(conn net.Conn) Client {
//	return Client{Conn: conn}
//}

// 不自动关闭连接，由调用使用方主动关闭连接
func (c *Client) Close() {
	//select {
	//case <- c.close:
	//	close(c.close)
	//	//c.Conn.Close()
	//	fmt.Println("正常关闭")
	//case <- time.After(constant.TIME_OUT):
	//	fmt.Println("超时关闭")
	//	close(c.close)
	//	//c.Conn.Close()
	//}
}

// 客户端调用
// serviceName：调用的服务名
// methodName：调用的方法名
// args：调用的参数
// reply: 返回值
func (c *Client) Call(serviceName, methodName string, args []interface{}, reply ...interface{}) error {
	// 用于等待响应结果
	w := make(chan uint8)
	// 关闭连接
	defer c.Close()
	// 组装调用信息，并序列化，写入
	go c.request(serviceName, methodName, args)
	// 返回结果
	go c.response(w, reply...)
	<-w
	return nil
}

func (c *Client) NewConn(newOrOld bool) net.Conn {
	var addr string
	if newOrOld {
		addr = c.e.Resolve(c.b)
		c.addr = addr
	}
	conn, err := loadConn(addr)
	if err != nil {
		c.mux.Lock()
		defer c.mux.Unlock()
		conn = storeConn(addr)
	}
	return conn
}

func (c *Client) request(serviceName, methodName string, args []interface{}) error {
	// 组装调用信息，并序列化，写入
	call := &client.Header{
		ServiceName: serviceName,
		MethodName:  methodName,
		Args:        args,
	}

	data, err := json.Marshal(call)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 必须以换行符作为结尾（与服务端定好的协议）
	data = append(data, []byte("\n")...)
	err = c.write(data, 2)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (c *Client) write(data []byte, tryTimes int) error {
	conn := c.NewConn(true)
	_, err := conn.Write(data)
	if err != nil {
		// 连接出错之后，再尝试连接一次
		delConn(c.addr)
		tryTimes--
		if tryTimes > 0 {
			c.write(data, tryTimes)
		}
		return err
	}
	c.conn <- conn
	return nil
}

// 读取服务端返回的结果
func (c *Client) response(w chan uint8, reply ...interface{}) error {
	conn := <-c.conn
	//
	bufferReader := bufio.NewReader(conn)
	for {
		context, err := bufferReader.ReadString(constant.END_SIGN)
		if err != nil {
			fmt.Println(err)
			delConn(c.addr)
			w <- 1
			return err
		}
		if len(context) == 0 {
			continue
		}
		var rsp client.Response
		err = json.Unmarshal([]byte(context), &rsp)
		if err != nil {
			fmt.Println(err)
			w <- 1
			return err
		}
		if rsp.Code != constant.SUCCESS {
			w <- 1
			return fmt.Errorf("远程调用错误码：【%s】，错误信息：【%s】", rsp.Code, rsp.Message)
		}
		data, err := json.Marshal(rsp.Data)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = json.Unmarshal(data, &reply)
		if err != nil {
			fmt.Println(err)
			w <- 1
			return err
		}
		w <- 1
		// 一次请求的响应完成后，必须返回。否则会无限的循环
		return nil
	}
	return nil
}
