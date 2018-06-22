package balancer

import "github.com/coreos/etcd/clientv3"


const(
	BALANCE_ROUND_ROBIN = 1
)
// 负载均衡接口（根据不同的需要，后续添加）
type Balancer interface {
	Resovle(rsp *clientv3.GetResponse) string
}