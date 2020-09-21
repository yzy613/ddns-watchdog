```
     _______________   _______________   ______      ___   _______________
    /  _______    \ \ /  _______    \ \ /  _   |    /  /| /  ____________/|
   /  /|______\   / //  /|______\   / //  / |  |   /  / //  |____________|/
  /  / /      /  / //  / /      /  / //  / /|  |  /  / /|\__________  \ \
 /  /_/______/  / //  /_/______/  / //  / / |  |_/  / / _|_________/  / /
/______________/ //______________/ //__/ /  |______/ //______________/ /
|______________|/ |______________|/ |__|/   |______|/ |______________|/

```
# Watchdog-ddns
[![Language](https://img.shields.io/badge/language-Go-00acd7)](https://golang.org)
[![DeepSource](https://static.deepsource.io/deepsource-badge-light-mini.svg)](https://deepsource.io/gh/yzy613/watchdog-ddns/?ref=repository-badge)
[![Build Status](https://travis-ci.com/yzy613/watchdog-ddns.svg?branch=master)](https://travis-ci.com/yzy613/watchdog-ddns)
[![Releases](https://img.shields.io/github/v/release/yzy613/watchdog-ddns)](https://github.com/yzy613/watchdog-ddns/releases)
[![Downloads](https://img.shields.io/github/downloads/yzy613/watchdog-ddns/total)](https://github.com/yzy613/watchdog-ddns/releases)
[![ClickDownload](https://img.shields.io/badge/%E7%82%B9%E5%87%BB-%E4%B8%8B%E8%BD%BD-brightgreen)](https://github.com/yzy613/watchdog-ddns/releases)

开箱即用的可常驻 Dynamic DNS 客户端，现已支持 DNSPod AliDNS(阿里云 DNS) Cloudflare，支持使用网卡 IP

## 准备工作
- 拥有（动态）公网 IPv6 地址或（动态）公网 IPv4 地址

- 拥有一个属于你的域名且有权操作（若没有，可以前往 DNSPod 阿里云或其他提供商购买域名）

- 域名的 DNS 服务器指向 [支持的服务商](https://github.com/yzy613/watchdog-ddns#%E6%94%AF%E6%8C%81%E7%9A%84%E6%9C%8D%E5%8A%A1%E5%95%86)

## 客户端
- 自动将域名解析到动态 IP（支持 IPv6）

### 用法
- `./watchdog-ddns-client -init string` 有选择地初始化配置文件，可以组合使用 (例 01)
    ```
    0 -> client.json
    1 -> dnspod.json
    2 -> alidns.json
    3 -> cloudflare.json
    ```

- `./watchdog-ddns-client` 直接运行

- `./watchdog-ddns-client -conf_path string` 指定配置文件路径 (最好是绝对路径)(路径有空格请放在双引号中间)

- `./watchdog-ddns-client -install` 安装服务

- `./watchdog-ddns-client -uninstall` 卸载服务

- `./watchdog-ddns-client -f` 强制检查 DNS 解析记录

- `./watchdog-ddns-client -version` 查看当前版本并检查更新

### 第一次使用？
- 前往 [releases](https://github.com/yzy613/watchdog-ddns/releases) 选择下载符合自己系统的压缩包，解压得到二进制文件

- 注意：Windows 的记事本保存的文件编码为 UTF-8 with BOM，需要使用第三方编辑器手动重新编码为 UTF-8，否则将会出现乱码

- 在 Linux 上不要忘记程序需要执行权限 `chmod 755 watchdog-ddns-client`

- 使用 `./watchdog-ddns-client -init string` 初始化配置文件（在 Windows 上使用 [watchdog-ddns-client-startup-script.bat](https://github.com/yzy613/watchdog-ddns/blob/master/watchdog-ddns-client-startup-script.bat) 一气呵成）

- 默认获取的是 IPv4，若需 IPv6，请在 `./conf/client.json` 修改 `api_url` 为 `https://yzyweb.cn/watchdog-ddns6`（此方法会优先使用临时 IPv6 地址）

- 若需使用网卡的 IP 地址，请在 `./conf/client.json` 修改 `enable_network_card` 为 `true` 并运行一次程序自动获取网卡信息，从 `./conf/network_card.json` 里面选择一个网卡填入 `./conf/client.json` 的 `network_card`

- 按照 [支持的服务商](https://github.com/yzy613/watchdog-ddns#%E6%94%AF%E6%8C%81%E7%9A%84%E6%9C%8D%E5%8A%A1%E5%95%86) 进行配置

- 若需配置不同域名的 watchdog-ddns，可以结合 `-conf_path string` 启动参数配置多种配置文件（可搭配 `-init` 启动参数初始化配置文件）

- 如果域名解析更新成功，那么程序工作正常，可以在 `./conf/client.json` 启用 `check_cycle` 进行定期检查（单位：分钟）（默认为 0，意为不启用定期检查）

- Enjoy it!（喜欢可以点一个 star 哦）

### 可选操作
- 在 Linux 上
1. 使用 `./watchdog-ddns-client -install` 安装服务，就可以使用 `systemctl` 管理 watchdog-ddns-client 服务了

- 在 Windows 上
1. [watchdog-ddns-client-startup-script.bat](https://github.com/yzy613/watchdog-ddns/blob/master/watchdog-ddns-client-startup-script.bat) 直接运行程序并回显程序返回的信息（需与 watchdog-ddns-client.exe 同一文件夹）
2. [watchdog-ddns-client-nohup.vbs](https://github.com/yzy613/watchdog-ddns/blob/master/watchdog-ddns-client-nohup.vbs) 不弹出运行窗口，在后台静默运行（需与 watchdog-ddns-client.exe 同一文件夹）
3. `Win`+`R` 后键入 `shell:startup` 会打开开机启动文件夹，将快捷方式粘贴在此处，即可进行开机启动（或把 watchdog-ddns-client-nohup.vbs 的快捷方式粘贴在此处，进行开机启动后台静默运行）

### 支持的服务商
#### DNSPod
- 请在 `./conf/client.json` 修改 `dnspod` 为 `true`

- 打开配置文件 `./conf/dnspod.json` 填入你的 `id, token, domain, sub_domain` 并重新启动

- 支持同一个域名的多个子域名同时更新解析
    ```json
    {
        "id": "",
        "token": "",
        "domain": "example.com",
        "sub_domain": [
            "@",
            "wwww"
        ]
    }
    ```

#### AliDNS (阿里云 DNS)
- 请在 `./conf/client.json` 修改 `alidns` 为 `true`

- 打开配置文件 `./conf/alidns.json` 填入你的 `accesskey_id, accesskey_secret, domain, sub_domain` 并重新启动

- 支持同一个域名的多个子域名同时更新解析
    ```json
    {
        "accesskey_id": "",
        "accesskey_secret": "",
        "domain": "example.com",
        "sub_domain": [
            "@",
            "wwww"
        ]
    }
    ```

#### Cloudflare
- 请在 `./conf/client.json` 修改 `cloudflare` 为 `true`

- 打开配置文件 `./conf/cloudflare.json` 填入你的 `email, api_key, zone_id, domain` 并重新启动

- 支持多个域名同时更新解析
    ```json
    {
        "email": "",
        "api_key": "",
        "zone_id": "",
        "domain": [
            "example.com",
            "www.example.com"
        ]
    }
    ```

#### 没有找到你的域名解析服务商？
- 请在 [Issues](https://github.com/yzy613/watchdog-ddns/issues) 提出请求或者 pr（感激不尽）

## 服务端 (普通用户不会用到，请略过)
- 返回 Json 格式的客户端 IP（支持 IPv6）

### 服务端 用法
- `./watchdog-ddns-server -install` 安装服务 (已经包含 `-init` 启动参数)

- `./watchdog-ddns-server -conf_path string` 指定配置文件路径 (最好是绝对路径)(路径有空格请放在双引号中间)

- `./watchdog-ddns-server -init` 初始化配置文件

- `systemctl start watchdog-ddns-server` 启动服务

- `./watchdog-ddns-server -uninstall` 卸载服务

- `./watchdog-ddns-server -version` 查看当前版本并检查更新

## 安装 
### Arch Linux
- 基于 Arch Linux 的发行版, 可以通过 [AUR](https://aur.archlinux.org/packages/watchdog-ddns-client) 安装
```bash
yay -S watchdog-ddns          // yay
pikaur -S watchdog-ddns       // pikaur
```

## 文档资源
> DNSPod API [https://www.dnspod.cn/docs/index.html](https://www.dnspod.cn/docs/index.html)

> Aliyun SDK [GitHub](https://github.com/aliyun/alibaba-cloud-sdk-go) or [https://help.aliyun.com/product/29697.html](https://help.aliyun.com/product/29697.html)

> Cloudflare API [https://api.cloudflare.com/#dns-records-for-a-zone-properties](https://api.cloudflare.com/#dns-records-for-a-zone-properties)
