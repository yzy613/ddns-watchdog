```
     _______________   _______________   ______      ___   _______________
    /  _______    \ \ /  _______    \ \ /  _   |    /  /| /  ____________/|
   /  /|______\   / //  /|______\   / //  / |  |   /  / //  |____________|/
  /  / /      /  / //  / /      /  / //  / /|  |  /  / /|\__________  \ \
 /  /_/______/  / //  /_/______/  / //  / / |  |_/  / / _|_________/  / /
/______________/ //______________/ //__/ /  |______/ //______________/ /
|______________|/ |______________|/ |__|/   |______|/ |______________|/

```
# 动态域名解析
[![language](https://img.shields.io/badge/language-Go-00acd7)](https://golang.org)
![build](https://travis-ci.com/yzy613/ddns.svg?branch=master)
[![release](https://img.shields.io/github/v/release/yzy613/ddns)](https://github.com/yzy613/ddns/releases)
[![downloads](https://img.shields.io/github/downloads/yzy613/ddns/total)](https://github.com/yzy613/ddns/releases)

开箱即用的 Dynamic DNS 客户端，现已支持 DNSPod 阿里云

## 准备工作
- 拥有 IPv6 地址或动态公网 IPv4 地址

- 拥有一个属于你的域名并有权操作（若没有，可以前往 DNSPod、阿里云或其他提供商购买域名）

- 域名的 DNS 服务器指向 [支持的服务商](https://github.com/yzy613/ddns#%E6%94%AF%E6%8C%81%E7%9A%84%E6%9C%8D%E5%8A%A1%E5%95%86)

## 客户端
- 自动将域名解析到动态 IP（支持 IPv6）

### 用法
- `./ddns-client` 直接运行

- `./ddns-client -mt` 显示更多的提示

- `./ddns-client -f` 强制检查 DNS 解析记录

- `./ddns-client -version` 查看当前版本

### 第一次使用？
- 不要忘记程序需要执行权限 `chmod 774 ddns-client`

- 初始化过程中需要手动重启多次，建议您使用终端或控制台运行程序并留意返回的信息

- 默认获取的是 IPv4，若需 IPv6，请在 `./conf/client.json` 修改 `web_addr` 为 `https://yzyweb.cn/ddns6`

- 按照 [支持的服务商](https://github.com/yzy613/ddns#%E6%94%AF%E6%8C%81%E7%9A%84%E6%9C%8D%E5%8A%A1%E5%95%86) 进行配置

- 如果程序无提示结束且域名解析未更新，那么请添加 `-mt` 启动参数以显示更多的提示

- 如果程序无提示结束且域名解析更新成功，那么程序工作正常，可以添加至计划任务定时执行

- Cron 命令 `*/10 * * * * cd /opt/ddns/ && ./ddns-client` 这里举例的位置是 `/opt/ddns/` 请根据自己的情况更改

- Enjoy it!

### 支持的服务商

#### DNSPod
- 请打开配置文件 `./conf/dnspod.json` 填入你的 `id, token, domain, sub_domain` 并重新启动

- 如果没有 `./conf/dnspod.json` 配置文件，请注意是否在 `./conf/client.json` 启用 `dnspod`

#### 阿里云（万网）
- 请打开配置文件 `./conf/aliyun.json` 填入你的 `accesskey_id, accesskey_secret, domain, sub_domain` 并重新启动

- 如果没有 `./conf/aliyun.json` 配置文件，请注意是否在 `./conf/client.json` 启用 `aliyun`

#### 没有找到你的域名服务商？
- 请在 [issues](https://github.com/yzy613/ddns/issues) 提出请求或者 pr 做出贡献 (感激不尽)

## 服务端 (一般情况下，不使用。请略过)
- 返回 Json 格式的客户端 IP（支持 IPv6）

### 服务端 用法
- `./ddns-server -install` 安装服务

- 使用 `systemctl start ddns-server` 启动

- `./ddns-server -uninstall` 卸载服务

- `./ddns-server -version` 查看当前版本

## 文档资源
- DNSPod API (https://www.dnspod.cn/docs/index.html)

- Aliyun SDK [GitHub](https://github.com/aliyun/alibaba-cloud-sdk-go) or (https://help.aliyun.com/product/29697.html)