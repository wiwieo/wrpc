package random

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"math/rand"
	"time"
)

type Random struct {
}

func (rb *Random) Resovle(rsp *clientv3.GetResponse) string {
	kvs := rsp.Kvs
	if len(kvs) == 0 {
		return ""
	}
	var addrs []string = make([]string, 0, len(kvs))
	for _, v := range kvs {
		addrs = append(addrs, string(v.Value))
	}
	fmt.Println(addrs)
	return addrs[rand.Int63n(time.Now().Unix())%int64(len(addrs))]
}
