package main

import (
	"net"
	"fmt"
	"wrpc/werror"
	"wrpc/client"
	"flag"
	"wrpc/etcdv3"
	"strings"
	"context"
	"wrpc/etcdv3/balancer/random"
)

type Param struct {
	One, Two int
	SonParam
}

type SonParam struct {
	SOne, STwo int
}
// 并发访问远程方法
//func main() {
//	var wg sync.WaitGroup
//	fmt.Println("开始时间：%v", time.Now())
//	for i:=0;i<10000;i++{
//		wg.Add(1)
//		go func(wg *sync.WaitGroup, i int){
//			defer wg.Done()
//			conn, err := net.Dial("tcp", "localhost:3333")
//			werror.CheckErr(err)
//			defer func() {
//				fmt.Println("中止")
//				conn.Close()
//			}()
//			c := client.NewClient(conn)
//			var rply int
//			err = c.Call("GateWay", "Add3", []interface{}{Param{One:1, Two:i}}, &rply)
//			if err != nil {
//				fmt.Println(err)
//			}
//			fmt.Println(fmt.Sprintf("第%d个add3的结果：%v", i, rply))
//		}(&wg, i)
//	}
//	wg.Wait()
//	fmt.Println("结束时间：%v", time.Now())
//}

var(
	HOST = flag.String("host", "localhost:33333", "添加环境变量，如 [-host=localhost:33333]")
	USE_ETCD = flag.Bool("use_etcd", true, "添加环境变量，如 [-use_etcd=true|false]")
	SERVERNAME = flag.String("server_name", "wiwieo/wrpc", "添加环境变量，如 [-server_name=wiwieo/wrpc]")
	ENDPOINTS = flag.String("endpoints", "http://localhost:2379", "添加环境变量，多个以逗号分隔 如 [-endpoints=http://localhost:2379,http://localhost:2380]")
)

func init() {
	flag.Parse()
}

func main() {
	var addr = HOST
	if *USE_ETCD {
		// 服务注册进etcd中
		e, err := etcdv3.NewEtcdv3(&etcdv3.ServerInfo{
			Name: *SERVERNAME,
			Endpoints: strings.Split(*ENDPOINTS, ","),
			TTL: 10,
			Ctx: context.Background(),
		})
		werror.CheckError(err)
		b := &random.Random{}
		*addr = e.Resolve(b)
		go e.Watch()
	}

	// 连接服务器
	conn, err := net.Dial("tcp", *addr)
	werror.CheckError(err)
	defer func() {
		conn.Close()
	}()
	c := client.NewClient(conn)

	// 调用远程方法
	var rply int
	err = c.Call("GateWay", "Add3", []interface{}{Param{One:1, Two:2}}, &rply)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("add3的结果：%v", rply))
}


