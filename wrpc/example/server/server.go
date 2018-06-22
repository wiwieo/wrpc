package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"
	"wrpc/etcdv3"
	"wrpc/server"
)

type header struct {
	ServiceName string
	MethodName  string
	Args        []interface{}
	Reply       chan []reflect.Value
}

var (
	USE_ETCD   = flag.Bool("use_etcd", true, "添加环境变量，如 [-use_etcd=true|false]")
	PORT       = flag.String("port", "33333", "lisntern port")
	SERVERNAME = flag.String("server_name", "wiwieo/wrpc", "添加环境变量，如 [-server_name=wiwieo/wrpc]")
	ENDPOINTS  = flag.String("endpoints", "http://localhost:2379", "添加环境变量，多个以逗号分隔 如 [-endpoints=http://localhost:2379,http://localhost:2380]")
)

func init() {
	// 一开始没注意忘了写这行，总是取不到配置的值，备注不要忘了
	flag.Parse()
}

func main() {
	// 是否使用ETCD作服务注册发现
	if *USE_ETCD {
		// 服务注册进etcd中
		e := etcdv3.Register(&etcdv3.ServerInfo{
			Host:      "localhost",
			Port:      *PORT,
			Name:      *SERVERNAME,
			Endpoints: strings.Split(*ENDPOINTS, ","),
			TTL:       10,
			Ctx:       context.Background(),
		})
		fmt.Println(*PORT)
		// 当服务KILL之后，使用信号来作通知
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
		go func() {
			s := <-ch
			log.Printf("receive signal '%v'", s)
			e.UnRegisterServer()
			os.Exit(1)
		}()
	}

	// 开启RPC服务监听
	server.Register(&GateWay{}, "")
	server.StartServerForTCP(fmt.Sprintf("localhost:%s", *PORT))
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
	time.Sleep(3 * time.Second)
	return p.One + p.Two
}

func (gw GateWay) Add4() int {
	return 1 + 2
}
