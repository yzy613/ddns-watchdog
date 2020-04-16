# 动态域名解析
![](https://travis-ci.com/yzy613/ddns.svg?branch=master)


DDNS 的服务端和客户端

## 客户端 用法
- `直接运行`

- `./ddns-client -f` 强制刷新 DNS 解析记录

- `加入计划任务定时执行`

## 服务端 用法
- `./ddns-server -install` 安装服务

- 用 `systemctl start ddns-server` 启动