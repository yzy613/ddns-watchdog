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
[![DeepSource](https://static.deepsource.io/deepsource-badge-light-mini.svg)](https://deepsource.io/gh/yzy613/ddns/?ref=repository-badge)
![build](https://travis-ci.com/yzy613/ddns.svg?branch=master)
[![release](https://img.shields.io/github/v/release/yzy613/ddns)](https://github.com/yzy613/ddns/releases)
[![downloads](https://img.shields.io/github/downloads/yzy613/ddns/total)](https://github.com/yzy613/ddns/releases)
[![click_download](https://img.shields.io/badge/%E7%82%B9%E5%87%BB-%E4%B8%8B%E8%BD%BD-brightgreen)](https://github.com/yzy613/ddns/releases)

开箱即用的 Dynamic DNS 客户端，现已支持 DNSPod 阿里云 Cloudflare，支持网卡 IP

## 准备工作
- 拥有（动态）公网 IPv6 地址或（动态）公网 IPv4 地址

- 拥有一个属于你的域名且有权操作（若没有，可以前往 DNSPod、阿里云或其他提供商购买域名）

- 域名的 DNS 服务器指向 [支持的服务商](https://github.com/yzy613/ddns#%E6%94%AF%E6%8C%81%E7%9A%84%E6%9C%8D%E5%8A%A1%E5%95%86)

## 客户端
- 自动将域名解析到动态 IP（支持 IPv6）

### 用法
- `./ddns-client -init` 初始化配置文件

- `./ddns-client` 直接运行

- `./ddns-client -conf_path string` 手动设置配置文件路径（绝对路径）（有空格用双引号）

- `./ddns-client -f` 强制检查 DNS 解析记录

- `./ddns-client -version` 查看当前版本并检查更新

### 第一次使用？
- 注意：Windows 的记事本保存的文件编码为 UTF-8 with BOM，需要使用第三方编辑器手动重新编码为 UTF-8，否则会产生乱码

- 不要忘记程序需要执行权限 `chmod 755 ddns-client`

- 使用 `./ddns-client -init` 初始化配置文件

- 默认获取的是 IPv4，若需 IPv6，请在 `./conf/client.json` 修改 `api_url` 为 `https://yzyweb.cn/ddns6`

- 按照 [支持的服务商](https://github.com/yzy613/ddns#%E6%94%AF%E6%8C%81%E7%9A%84%E6%9C%8D%E5%8A%A1%E5%95%86) 进行配置

- 若需使用网卡的 IP 地址，请在 `./conf/client.json` 修改 `enable_network_card` 为 `true` 并运行一次程序自动获取网卡信息，从 `./conf/network_card.json` 里面选择一个网卡填入 `./conf/client.json` 的 `network_card`

- 若需配置多个 IP 的 ddns，可以结合 `-conf_path string` 启动参数配置多种配置文件（可搭配 `-init` 启动参数初始化配置文件）

- 如果程序无提示结束且域名解析更新成功，那么程序工作正常，可以在 `./conf/client.json` 启用 `check_cycle` 进行定期检查（单位：分钟）（默认为 0，意为不启用定期检查）

- Enjoy it!

### 支持的服务商

#### DNSPod
- 请打开配置文件 `./conf/dnspod.json` 填入你的 `id, token, domain, sub_domain` 并重新启动

- 如果没有生效，请注意是否在 `./conf/client.json` 启用 `dnspod`

#### 阿里云（万网）
- 请打开配置文件 `./conf/aliyun.json` 填入你的 `accesskey_id, accesskey_secret, domain, sub_domain` 并重新启动

- 如果没有生效，请注意是否在 `./conf/client.json` 启用 `aliyun`

### Cloudflare
- 请打开配置文件 `./conf/cloudflare.json` 填入你的 `email, api_key, zone_id, domain` 并重新启动

- 如果没有生效，请注意是否在 `./conf/client.json` 启用 `cloudflare`

#### 没有找到你的域名解析服务商？
- 请在 [Issues](https://github.com/yzy613/ddns/issues) 提出请求或者 pr (感激不尽)

## 服务端 (普通用户不会用到。请略过)
- 返回 Json 格式的客户端 IP（支持 IPv6）

### 服务端 用法
- `./ddns-server -install` 安装服务

- `systemctl start ddns-server` 启动服务

- `./ddns-server -uninstall` 卸载服务

- `./ddns-server -version` 查看当前版本

## 文档资源
> DNSPod API (https://www.dnspod.cn/docs/index.html)

> Aliyun SDK [GitHub](https://github.com/aliyun/alibaba-cloud-sdk-go) or (https://help.aliyun.com/product/29697.html)

> Cloudflare API (https://api.cloudflare.com/#dns-records-for-a-zone-properties)