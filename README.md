# PacketVerify
Per-Packet capability verifiaction.

### Build

```
# build
go build

#
./packet-verify
```


Kill process in port
```
writecode@debian:~/dev/packet-verify$ netstat -tupln | grep 8080
(Not all processes could be identified, non-owned process info
 will not be shown, you would have to be root to see it all.)
tcp6       0      0 :::8080                 :::*                    LISTEN      10009/./packet-veri 
writecode@debian:~/dev/packet-verify$ kill -9 10009
```

### 运行方式
目前已经实现echo server，可修改 tcp_client变成并发发多个请求，来确定请求延时。  
已实现：  
- 包转发和签名验证，现在tcp每个包都会验证  

未实现：
- Diffie Hellman交换密钥，目前使用是固定密钥  
- 丢包同步 （目前写了一部分，但是没完全实现）  
- proxy 处理多client，目前支持多client通过proxy建立多个通道去连server，需要完善测试代码

```bash
./start_server.sh

# 在client 这里输入数据，服务器会echo回来
./start_client.sh 
TCPClient Connect to: 127.0.0.1:8082
adad
Send Pkt Len:5
ECHO:adad
```

性能评估方法：  
```bash
# 将重要时间戳加log，utils.Statistic()
# 利用python解析log，获取时间戳
# 可利用packet size来识别同一个请求链路（tcp client发送不同size的包）
# LocalProxy log
2021/04/27 13:20:59.167119 [STATISTIC] [12] LocalProxy forward to server: 16 sign: 39326
2021/04/27 13:21:00.648403 [STATISTIC] [4] LocalProxy forward to server: 8 sign: 648
2021/04/27 13:21:00.648519 [STATISTIC] [12] LocalProxy forward to server: 16 sign: 62865

# RemoteProxy log
2021/04/27 13:20:59.167472 [STATISTIC] [12] RemoteProxy forwarding to server: 12 sign: 39326
2021/04/27 13:20:59.167872 [STATISTIC] RemoteProxy forwarding to local: 13 sign: 57187
2021/04/27 13:21:00.648430 RemoteProxy recv from local: 8 4 648
2021/04/27 13:21:00.648563 [STATISTIC] [4] RemoteProxy forwarding to server: 4 sign: 648
2021/04/27 13:21:00.648631 RemoteProxy recv from local: 16 12 62865
2021/04/27 13:21:00.648690 [STATISTIC] [12] RemoteProxy forwarding to server: 12 sign: 62865
2021/04/27 13:21:00.648854 [STATISTIC] RemoteProxy forwarding to local: 17 sign: 80352
```

### Design

1. Connection Auth模块  
一个单独的服务用来进行设备和应用认证，完成认证之后eBPF firewall才会去打开端口。  
在验证完应用身份和设备Capability之后，会生成一个Token。后面的包和连接基于这个Packet进行认证。  

2. TCP/UDP Socket Hook  
方案一：系统层面替换TCP/UDP API，用真实IoT应用去评估  
直接修改或者Hook Socket接口，例如修改Kernel。或者利用netfilterqueue去修改TCP/UDP包  
或者使用LD_PRELOAD加载自定义动态库替换掉socket接口：https://my.oschina.net/xieyunzi/blog/669349  
通过eBPF对每一个包的签名进行检查。  

方案二：纯Proxy，无需eBPF，利用Proxy检查权限并将消息转发给新应用。  
https://github.com/snail007/goproxy/blob/master/utils/serve-channel.go  
当前实现方案：基于方案二。  

### 测试
使用nc客户端，或者直接使用go客户端(go run main.go --mode=TCPClient)  
nc 127.0.0.1 8080  

UDP:  
nc -u 127.0.0.1 8080  

UDP抓包：
sudo tcpdump udp -i any  

### TODO  
1. 基于Proxy来实现验证  
2. 搞清楚能否单UDP包添加签名  

3. 实现参考  
https://github.com/snail007/goproxy  
https://grpc.io/docs/languages/go/quickstart/  


Simple Proxy:  
https://github.com/arkadijs/goproxy  
https://gist.github.com/mike-zhang/3853251  

4. 简单演示

client -> proxy(8000) -> server(8080)  
```  
go build
./packet-verify --mode=TCPServer
./packet-verify --mode=EbpfProxy
./packet-verify --mode=TCPClient
```  

client -> local_proxy -> | -> remote_proxy -> server  