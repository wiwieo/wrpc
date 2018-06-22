package etcdv3

import (
	"github.com/coreos/etcd/clientv3"
	"wrpc/werror"
	"wrpc/etcdv3/balancer"
)

// 从etcd中获取对应的数据，并解析
func (e *Etcdv3) Resolve(b balancer.Balancer) string{
	rsp, err := e.client.Get(e.s.Ctx, e.s.GetPrefixKey(), clientv3.WithPrefix())
	werror.CheckError(err)
	return b.Resovle(rsp)
}