package etcdv3

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"net"
	"wrpc/werror"
)

type ServerInfo struct {
	Name      string          // 服务名 如：wrpc
	Host      string          // 服务主机 如：127.0.0.1
	Port      string          // 服务端口号 如：3333
	Endpoints []string        // etcd3的连接信息 如：http://localhost:2379
	TTL       int64           // 最后存活时间
	Ctx       context.Context // 上下文
}

func (s *ServerInfo) GetPrefixKey() string {
	return fmt.Sprintf("/%s/", s.Name)
}

func (s *ServerInfo) GetFullKey() string {
	return fmt.Sprintf("/%s/%s", s.Name, net.JoinHostPort(s.Host, s.Port))
}

func (s *ServerInfo) GetHostPort() string {
	return net.JoinHostPort(s.Host, s.Port)
}

type Etcdv3 struct {
	client *clientv3.Client
	s      *ServerInfo
}

func NewEtcdv3(s *ServerInfo) (*Etcdv3, error) {
	c, err  := clientv3.New(clientv3.Config{Endpoints: s.Endpoints})
	if err != nil{
		return nil, err
	}
	return &Etcdv3{
		client: c,
		s:      s,
	}, nil
}

// 服务注册
func Register(s *ServerInfo) *Etcdv3 {
	e, err := NewEtcdv3(s)
	werror.CheckError(err)
	e.registerServer()
	return e
}

func (e *Etcdv3) registerServer() {
	// 创建租约
	resp, err := e.client.Grant(e.s.Ctx, e.s.TTL)
	if err != nil {
		werror.CheckError(err)
	}

	// 存放进etcd中
	_, err = e.client.Put(e.s.Ctx, e.s.GetFullKey(), e.s.GetHostPort(), clientv3.WithLease(resp.ID))
	if err != nil {
		werror.CheckError(err)
	}

	// 确保租约一直有效
	if _, err := e.client.KeepAlive(e.s.Ctx, resp.ID); err != nil {
		werror.CheckError(err)
	}

	// wait deregister then delete
	//go func() {
	//	<-Deregister
	//	client.Delete(ctx, serviceKey)
	//	Deregister <- struct{}{}
	//}()

	//return nil
}

// 服务注销
func (e *Etcdv3) UnRegisterServer() {
	e.client.Delete(e.s.Ctx, e.s.GetFullKey())
}
