# [GoSocksProxy](https://github.com/Yu2erer/GoSocksProxy)
一个 Golang 所写的网络混淆代理 仅用来 学习 SOCKS5 代理. 实现过程可见Blog [实现 SOCKS5 网络混淆代理](https://www.yuerer.com/实现-SOCKS5-网络混淆代理/)

![GoSocksProxy](https://www.yuerer.com/images/GoSocksProxy.jpg)
## 使用方法

* 修改 server/server.go 的密码(默认为 "haohaio!") 并将其编译 在墙外服务器运行即可
* 修改 local/local.go 的墙外服务器IP地址 和 密码(默认为 "haohaio!") 并将其编译 运行于本地即可
* 默认 本地监听地址为 `127.0.0.1:1234`
* 需要将本地的代理服务器地址设置为以上地址 协议为 SOCKS5

## 注意

* 默认只支持 `DES` 加密. 因为本项目只是为了学习 SOCKS5 协议.
* server/server.go 写的很乱 请小心驾驶

## TODO(不存在的ToDo)

* 加入更多加密算法
* 重写 SOCKS5 协议的 握手 和 获取请求部分