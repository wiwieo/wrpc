# wrpc
为了加深RPC的理解，尝试写的简易RPC框架，用ETCD3作为服务发现，支持负载均衡

# 相关技术
* 1、使用`TCP` + `json`作为协议
 
* 2、远程调用时，使用反射进行实际的调用

* 3、用etcd3来作服务注册发现及注销

# 在写该项目时，遇到的问题

一、协议的选用，是`UDP`还是`TCP`，之后考虑到简单性选择了`TCP`。两者区别参考如下文章

* TCP和UDP的优缺点及区别：（https://www.cnblogs.com/xiaomayizoe/p/5258754.html）

二、方法调用问题

* 问题描述：GO中的反射无法直接根据`结构体名及方法名`直接映射成对应的实例并调用（java可以直接映射成对应的实例），如何将远程读取到的二进制数据转换成对应的方法并调用？
* 解决方案：服务启动时，将需要暴露的方法注册进绶存中。如：使用map来存放，将`服务名.方法名`作key，将调用的信息经过反射后作为value。如此，当接收到信息时，直接根据调用的方法来获取对应的实例并执行。还可以进行方法过虑，不会调用到未暴露的方法。

三、序列化问题。

* 问题描述：在服务端接收到请求之后，根据事前约定的协议进行反序列化。在具体调用方法时，如何将参数是自定义复杂结构实例化。
* 解决方案：此问题和上一个问题类似。直接根据获取到的实例信息，使用`reflect.New`来创建对应的实例，然后将对应的值赋其之即可。

四、负载均衡
* 问题描述：因为RPC使用的是长连接，那么负载均衡如何做到，一个请求到来，如何实时分发到多个服务器，而不是只发到一个已经连接到的服务器中？
* 解决方案：在项目启动时并不做实际连接，根据具体的负载方案（轮询，随机，盐值等）先确定使用哪一台服务，在实际调用时，再连接。且在连接处，使用了缓存，将每个连接缓存起来，根据连接的地址来获取对应的连接，如果获取到则直接使用。如果获取不到，则建立连接，并存放到缓存中。在服务断连之后，将连接从缓存中删除

五、服务断开连接通知
* 问题描述一：如何知道服务断开并告知客户端，将缓存中的连接删除
* 解决方案：使用信号，当服务断开后，使用SIGN通知到程序，将etcd中的key-value删除。

* 问题描述二：何时删除客户端维护的连接池（简便起见，使用的是方案二）
* 解决方案一：可以使用etcd的watch功能，当key-value发生变化时，可以通过程序监听并做相应的操作
* 解决方案二：该问题的影响在于服务断开又连接时，客户端会使用到旧连接，会导致读写数据失败。此处是在读写数据失败时，将连接删除，并重新连接。