package round_robin

import (
	"github.com/coreos/etcd/clientv3"
	"sync/atomic"
)

type RoundRobin struct {
	Count int64
}

func (rb *RoundRobin) Resovle(rsp *clientv3.GetResponse) string {
	kvs := rsp.Kvs
	var addrs []string = make([]string, 0, len(kvs))
	for _, v := range kvs {
		addrs = append(addrs, string(v.Value))
	}
	addr := addrs[rb.Count%int64(len(addrs))]
	rb.Count = atomic.AddInt64(&rb.Count, 1)
	return addr
}
