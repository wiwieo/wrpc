package etcdv3

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"fmt"
)

// 监听etcd值的动态
func (e *Etcdv3) Watch(){
	wc := e.client.Watch(e.s.Ctx, e.s.GetPrefixKey(), clientv3.WithPrefix())
	for wresp := range wc {
		for _, ev := range wresp.Events {
			fmt.Println(fmt.Sprintf("etcd3服务地址：%v", string(ev.Kv.Value)))
			switch ev.Type {
			case mvccpb.DELETE:
			}
		}
	}
}