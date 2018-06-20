package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"wrpc/constant"
	"wrpc/entity/client"
	"time"
)

type Client struct {
	Conn net.Conn
	close chan uint8
}

func NewClient(conn net.Conn) Client {
	return Client{Conn: conn}
}

func (c *Client) Close(){
	select {
	case <- c.close:
		c.Conn.Close()
		fmt.Println("正常关闭")
	case <- time.After(constant.TIME_OUT):
		fmt.Println("超时关闭")
		c.Conn.Close()
	}
}

// 客户端调用
// serviceName：调用的服务名
// methodName：调用的方法名
// args：调用的参数
// reply: 返回值
func (c *Client) Call(serviceName, methodName string, args []interface{}, reply ...interface{}) error {
	c.close = make(chan uint8)
	// 关闭连接
	defer c.Close()
	// 组装调用信息，并序列化，写入
	go c.request(serviceName, methodName, args)
	// 返回结果
	go c.response(reply...)
	return nil
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
		return err
	}
	// 必须以换行符作为结尾（与服务端定好的协议）
	data = append(data, []byte("\n")...)
	_, err = c.Conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// 读取服务端返回的结果
func (c *Client) response(reply ...interface{}) error {
	// 需要自旋来获取返回结果
	for {
		bufferReader := bufio.NewReader(c.Conn)
		context, _, err := bufferReader.ReadLine()
		if err != nil {
			return err
		}
		if len(context) == 0 {
			continue
		}
		var rsp client.Response
		err = json.Unmarshal(context, &rsp)
		if err != nil {
			return err
		}
		if rsp.Code != constant.SUCCESS {
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
			return err
		}
		c.close <- 1
		return nil
	}
	return nil
}
