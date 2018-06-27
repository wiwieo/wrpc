package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"wrpc/client"
	"wrpc/etcdv3/balancer"
)

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

type Example struct {
	C *client.Client
}

var (
	HOST       = flag.String("host", "localhost:33333", "添加环境变量，如 [-host=localhost:33333]")
	PORT       = flag.String("port", "8080", "添加环境变量，如 [-port=8080]")
	USE_ETCD   = flag.Bool("use_etcd", true, "添加环境变量，如 [-use_etcd=true|false]")
	SERVERNAME = flag.String("server_name", "wiwieo/wrpc", "添加环境变量，如 [-server_name=wiwieo/wrpc]")
	ENDPOINTS  = flag.String("endpoints", "http://localhost:2379", "添加环境变量，多个以逗号分隔 如 [-endpoints=http://localhost:2379,http://localhost:2380]")
)

func init() {
	flag.Parse()
}

func main() {
	c, err := client.DialTCPConnect(*USE_ETCD, *HOST, balancer.BALANCE_ROUND_ROBIN, *SERVERNAME, *ENDPOINTS)
	if err != nil {
		fmt.Println("系统连接报错")
		os.Exit(-1)
	}
	e := &Example{
		C: c,
	}
	e.StartWeb()
}

func (e *Example) StartWeb() {
	http.HandleFunc("/add", e.Add)

	http.ListenAndServe(fmt.Sprintf(":%s", *PORT), nil)
}
