# TunnHub - NetworkTunnel Hub

<br>

[中文文档](./tunnhub_cn.md) | [English](./tunnhub_en.md)

<br>

### 特性

--------

#### 支持的操作系统

已测试的操作系统：

- Windows 7/10/11
- CentOS 7.x
- Ubuntu 20.x

理论支持：

- Windows 7+
- 支持tun设备的Linux发行版

#### 支持的传输协议

TCP / KCP / WS / WSS

#### 支持的加密方式

AES256 / AES192 / AES128 / XOR / SM4 / TEA / XTEA / Salsa20 / Blowfish

### 更新

------

2022/05/27 @ 1.0.0.220527

- 验证服务器重构
- 证书管理
- 前端WS支持
- 网络暴露/导入流程调整
- 系统运行监控
- 网卡数据保存
- 远程重置/停止连接
- 静态地址分配

2022/05/10 @ 0.0.1

- 分离自项目 [Tunnel](https://gitee.com/jackrabbit872568318/tunnel)

### 编译

------

需要安装Go1.18.2或者更高版本 [下载](https://golang.google.cn/dl/) <br>
编译时需要GCC环境，Windows用户请安装 [mingw-w64](https://www.mingw-w64.org/)

准备

```shell
#拉取仓库
git clone https://gitee.com/jackrabbit872568318/tunn-hub.git

#进入目录
cd ./tunn-hub

#下载依赖
set GO111MODULE=on
go mod tidy

#进入cmd目录
cd cmd
```

编译

```shell
# @linux
go build -o tunnhub
```

或

```shell
# @windows
go build -o tunnhub.exe
```

### 使用

------

#### 客户端配置示例

[配置文件](../config/tunnhub_config_full.json)

说明

```shell
#通信IP地址，不指定地址使用(0.0.0.0)
global.address
#通信端口
global.port
#通信协议
global.protocol
#最大传输单元 (默认1400即可)
global.mtu
#客户端并行连接数
global.multi_connection

#路由名称
route.name
#在Hub中不会执行import操作，仅暴露网络
route.option
#网络CIDR，如：192.168.0.0/24，192.168.1.254/32
route.network

#虚拟网卡CIDR
device.cidr
#虚拟网卡DNS (仅Windows)
device.dns

#认证服务器地址，不指定地址使用(0.0.0.0)
auth.address
#认证服务器端口
auth.port

#数据加密处理模式，留空则不加密
data_process.encrypt

#证书文件路径
security.cert
#密钥文件路径
security.key

#控制台地址，不指定地址使用(0.0.0.0)
admin.address
#控制台端口
admin.port
#控制台Reporter端口(前端获取数据需要)
admin.reporter
#控制台是否开启Https访问
admin.https
#sqlite数据库文件
admin.db_file

#地址池动态分配起始地址
ip_pool.start
#地址池动态分配结束地址
ip_pool.end
#地址池网络，如：192.168.10.0/24
ip_pool.network
```

示例

```json
{
  "global": {
    "address": "127.0.0.1",
    "port": 10240,
    "protocol": "wss",
    "mtu": 1400,
    "multi_connection": 1
  },
  "route": [
    {
      "name": "test_route",
      "option": "export",
      "network": "10.10.10.10/32"
    }
  ],
  "device": {
    "cidr": "172.22.0.1/24",
    "dns": "223.5.5.5"
  },
  "auth": {
    "address": "0.0.0.0",
    "port": 10241
  },
  "data_process": {
    "encrypt": "XOR"
  },
  "security": {
    "cert": "./cert.pem",
    "key": "./key.pem"
  },
  "admin": {
    "address": "0.0.0.0",
    "port": 8888,
    "reporter": 8889,
    "https": true,
    "db_file": "./tunn_server.db"
  },
  "ip_pool": {
    "start": "172.22.0.11",
    "end": "172.22.0.100",
    "network": "172.22.0.0/24"
  }
}
```

#### 启动

! Windows需要以管理员模式启动 <br>
! Windows需要下载 [wintun](https://www.wintun.net/) 驱动并与可执行文件在同一目录下

启动参数

- -c 指定配置文件路径

示例：

```shell
# @linux
./tunnhub -c config.json
```

或

```shell
# @windows
tunnhub.exe -c config.json
```

启动成功时控制台如下

![img](./img/hub_startup.png)

使用https://your.hub.address:8888即可访问控制台

#### 第一次使用时

    控制台默认用户： admin
    控制台默认密码： P@ssw0rd

登录到控制台，进入证书管理页<br>
点击高级设置，将允许的IP或域名设置为你的服务器地址,设置证书过期时间并勾选覆盖配置
![img.png](img/cert_create.png)
如图设置客户端将可以使用域名"domain.tunnhub.com"和IP地址"192.168.0.85","123.123.123.123"连接， 并且证书会在2022-06-29 00:00:00过期

点击创建并重启服务器即可在页面下方下载证书，将证书提供给客户端


