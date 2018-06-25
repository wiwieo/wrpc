package balancer

import (
	"github.com/coreos/etcd/clientv3"
	"wrpc/etcdv3/balancer/random"
	"wrpc/etcdv3/balancer/roud_robin"
)

const (
	// 轮询
	BALANCE_ROUND_ROBIN = 1
	// 随机
	BALANCE_RANDOM      = 2
)

// 负载均衡接口（根据不同的需要，后续添加）
type Balancer interface {
	Resovle(rsp *clientv3.GetResponse) string
}

func GetBalance(balanceType int) Balancer{
	var b Balancer
	switch balanceType {
	case BALANCE_RANDOM:
		b = &random.Random{}
	case BALANCE_ROUND_ROBIN:
		b = &round_robin.RoundRobin{}
	}
	return b
}