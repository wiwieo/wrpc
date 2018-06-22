package round_robin

import (
	"github.com/coreos/etcd/clientv3"
)

type RoundRobin struct {
}

func (rb *RoundRobin) Resovle(rsp *clientv3.GetResponse) string {
	kvs := rsp.Kvs
	var addrs []string = make([]string, 0, len(kvs))
	for _, v:= range kvs{
		addrs = append(addrs, string(v.Value))
	}
	return addrs[0]
}
