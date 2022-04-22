```
     _______________   _______________   ______      ___   _______________
    /  _______    \ \ /  _______    \ \ /  _   |    /  /| /  ____________/|
   /  /|______\   / //  /|______\   / //  / |  |   /  / //  |____________|/
  /  / /      /  / //  / /      /  / //  / /|  |  /  / /|\__________  \ \
 /  /_/______/  / //  /_/______/  / //  / / |  |_/  / / _|_________/  / /
/______________/ //______________/ //__/ /  |______/ //______________/ /
|______________|/ |______________|/ |__|/   |______|/ |______________|/

```

# ddns-watchdog

[![Language](https://img.shields.io/badge/language-Go-00acd7)](https://golang.org)
[![DeepSource](https://static.deepsource.io/deepsource-badge-light-mini.svg)](https://deepsource.io/gh/yzy613/ddns-watchdog/?ref=repository-badge)
[![.github/workflows/go.yml](https://github.com/yzy613/ddns-watchdog/actions/workflows/go.yml/badge.svg)](https://github.com/yzy613/ddns-watchdog/actions/workflows/go.yml)
[![Releases](https://img.shields.io/github/v/release/yzy613/ddns-watchdog)](https://github.com/yzy613/ddns-watchdog/releases)
[![Downloads](https://img.shields.io/github/downloads/yzy613/ddns-watchdog/total)](https://github.com/yzy613/ddns-watchdog/releases)
[![ClickDownload](https://img.shields.io/badge/%E7%82%B9%E5%87%BB-%E4%B8%8B%E8%BD%BD-brightgreen)](https://github.com/yzy613/ddns-watchdog/releases)

现已支持 DNSPod AliDNS(阿里云 DNS) Cloudflare，支持 IPv4 IPv6 双栈，支持使用网卡 IP 地址

## 准备工作

1. 拥有 (公网动态) IPv6 地址或 (公网动态) IPv4 地址
2. 拥有一个能申请修改解析记录的 Token 的域名 (若没有，可以前往 DNSPod 阿里云或其他提供商购买域名)
3. 域名的 DNS 服务器指向 [支持的服务商](https://github.com/yzy613/ddns-watchdog#%E6%94%AF%E6%8C%81%E7%9A%84%E6%9C%8D%E5%8A%A1%E5%95%86)

## 客户端

### 用法

- `./ddns-watchdog-client -i 0123` 有选择地初始化配置文件并退出，可以组合使用

  此示例展示仅初始化客户端和 DNSPod 的配置文件

  ```bash
  ./ddns-watchdog-client -i 01
  ```

  代码参考表

  ```bash
  0 -> client.json
  1 -> dnspod.json
  2 -> alidns.json
  3 -> cloudflare.json
  ```
- `./ddns-watchdog-client` 使用默认配置文件目录 `./conf` 运行
- `./ddns-watchdog-client -n` 输出网卡信息并退出
- `./ddns-watchdog-client -c ./conf` 指定配置文件目录为 ./conf (目录有空格请放在双引号中间)
- `./ddns-watchdog-client -I` 安装服务并退出 (仅限有 systemd 的 Linux 使用)
- `./ddns-watchdog-client -U` 卸载服务并退出 (仅限有 systemd 的 Linux 使用)
- `./ddns-watchdog-client -f` 强制检查解析记录值
- `./ddns-watchdog-client -v` 查看当前版本并检查更新后退出

### 初始客户端配置文件

```json
{
    "api_url": {
        "ipv4": "https://yzyweb.cn/ddns-watchdog",
        "ipv6": "https://yzyweb.cn/ddns-watchdog6",
        "version": "https://yzyweb.cn/ddns-watchdog"
    },
    "enable": {
        "ipv4": false,
        "ipv6": false,
        "network_card": false
    },
    "network_card": {
        "ipv4": "",
        "ipv6": ""
    },
    "services": {
        "dnspod": false,
        "alidns": false,
        "cloudflare": false
    },
    "check_cycle_minutes": 0
}
```

### 第一次使用？

1. 前往 [releases](https://github.com/yzy613/ddns-watchdog/releases) 下载符合自己系统的压缩包，解压得到二进制文件
2. 注意：Windows 的记事本保存的文件编码为 UTF-8 with BOM，需要使用第三方编辑器手动重新编码为 UTF-8，否则将会出现乱码导致无法读取正确的配置
3. 在 Linux 上不要忘记程序需要执行权限 `chmod 755 ddns-watchdog-client`
4. 使用 `./ddns-watchdog-client -i 0123` 初始化配置文件 (在 Windows 上使用 [ddns-watchdog-client-startup-script.bat](https://github.com/yzy613/ddns-watchdog/blob/master/ddns-watchdog-client-startup-script.bat) 一气呵成)
5. 根据使用环境确定启用 (`enable`) IPv4 还是 IPv6 或是两者都启用
6. 若未启用网卡，默认使用 API 获取对应 IP 地址
7. 若需使用网卡的 IP 地址，请在 `./conf/client.json` 修改 `enable`->`network_card` 为 `true` 并运行一次程序自动获取网卡信息，从 `./conf/network_card.json` 里面选择网卡填入 `./conf/client.json` 的 `network_card`
8. 若 `./conf/client.json` 的 `network_card`->`ipv4` 或 `ipv6` 为空，对应 IP 地址将从 API 获取

   此示例展示 IPv4 从 API 获取，IPv6 从 example 网卡获取

   ```json
   {
       "enable": {
           "ipv4": true,
           "ipv6": true,
           "network_card": true
       },
       "network_card": {
           "ipv4": "",
           "ipv6": "example"
       }
   }
   ```
9. 按照 [支持的服务商](https://github.com/yzy613/ddns-watchdog#%E6%94%AF%E6%8C%81%E7%9A%84%E6%9C%8D%E5%8A%A1%E5%95%86) 进行配置
10. 若需配置不同域名的 ddns-watchdog，可以结合 `-c` 启动参数配置多种配置文件 (可搭配 `-i` 启动参数初始化配置文件)
11. 如果解析记录值更新成功，那么程序工作正常，可以在 `./conf/client.json` 启用 `check_cycle_minutes` 进行定期检查 (单位：分钟)(默认为 0，意为不启用定期检查)
12. 注意：ddns-watchdog 设计了 IP 地址本地比对机制，以防止频繁访问 API 导致封禁。若手动修改了解析记录值，会导致无法及时更新 (可搭配 `-f` 启动参数强制检查解析记录值以跳过本地比对机制)

    ***Enjoy it!（觉得好用可以点一个 star 噢）***

### 可选操作

- 在有 systemd (systemctl) 的 Linux 上

  1. 使用 `./ddns-watchdog-client -I` 安装服务，就可以使用 `systemctl` 管理 ddns-watchdog-client 服务了
- 在 Windows 上

  1. [ddns-watchdog-client-startup-script.bat](https://github.com/yzy613/ddns-watchdog/blob/master/scripts/run/ddns-watchdog-client-startup-script.bat) 一键运行程序并回显程序返回的信息 (需与 ddns-watchdog-client.exe 同一文件夹)
  2. [ddns-watchdog-client-nohup.vbs](https://github.com/yzy613/ddns-watchdog/blob/master/scripts/run/ddns-watchdog-client-nohup.vbs) 不弹出运行窗口，在后台静默运行 (需与 ddns-watchdog-client.exe 同一文件夹)
  3. `Win`+`R` 后键入 `shell:startup` 会打开开机启动文件夹，将快捷方式粘贴在此处，即可进行开机启动 (或把 ddns-watchdog-client-nohup.vbs 的快捷方式粘贴在此处，进行开机启动后台静默运行)

### 支持的服务商

#### DNSPod

- 请在 `./conf/client.json` 修改 `dnspod` 为 `true`
- 打开配置文件 `./conf/dnspod.json` 填入你的 `id, token, domain, sub_domain` 并重新启动
- 支持同一个域名的 A 和 AAAA 记录的子域名同时更新记录值

  初始 DNSPod 配置文件

  ```json
  {
      "id": "在 https://console.dnspod.cn/account/token/token 获取",
      "token": "在 https://console.dnspod.cn/account/token/token 获取",
      "domain": "example.com",
      "sub_domain": {
          "a": "A记录子域名",
          "aaaa": "AAAA记录子域名"
      }
  }
  ```

#### AliDNS (阿里云 DNS)

- 请在 `./conf/client.json` 修改 `alidns` 为 `true`
- 打开配置文件 `./conf/alidns.json` 填入你的 `accesskey_id, accesskey_secret, domain, sub_domain` 并重新启动
- 支持同一个域名的 A 和 AAAA 记录的子域名同时更新记录值

  初始 AliDNS 配置文件

  ```json
  {
      "accesskey_id": "在 https://ram.console.aliyun.com/users 获取",
      "accesskey_secret": "在 https://ram.console.aliyun.com/users 获取",
      "domain": "example.com",
      "sub_domain": {
          "a": "A记录子域名",
          "aaaa": "AAAA记录子域名"
      }
  }
  ```

#### Cloudflare

- 请在 `./conf/client.json` 修改 `cloudflare` 为 `true`
- 打开配置文件 `./conf/cloudflare.json` 填入你的 `zone_id, api_token, domain` 并重新启动
- 支持同一个域名的 A 和 AAAA 记录的子域名同时更新内容

  初始 Cloudflare 配置文件

  ```json
  {
      "zone_id": "在你域名页面的右下角有个区域 ID",
      "api_token": "在 https://dash.cloudflare.com/profile/api-tokens 获取",
      "domain": {
          "a": "A记录子域名.example.com",
          "aaaa": "AAAA记录子域名.example.com"
      }
  }
  ```

#### 没有找到你的域名解析服务商？

- 请在 [Issues](https://github.com/yzy613/ddns-watchdog/issues) 提出 Issue 或者在 [Pull requests](https://github.com/yzy613/ddns-watchdog/pulls) Pull request (感激不尽)

## 服务端 (普通用户不会用到，请略过)

返回 Json 格式的客户端 IP 地址 (支持 IPv4 IPv6 双栈)

### 服务端 用法

- `./ddns-watchdog-server -I` 安装服务并退出 (已经包含 `-i` 启动参数)
- `./ddns-watchdog-server -c ./conf` 指定配置文件目录为 ./conf (目录有空格请放在双引号中间)
- `./ddns-watchdog-server -i` 初始化配置文件并退出
- `systemctl start ddns-watchdog-server` 启动服务
- `./ddns-watchdog-server -U` 卸载服务并退出
- `./ddns-watchdog-server -v` 查看当前版本并检查更新后退出

## 安装

### Arch Linux

基于 Arch Linux 的发行版, 可以通过 [AUR](https://aur.archlinux.org/packages/ddns-watchdog-client) 安装

```bash
yay -S ddns-watchdog          // yay
pikaur -S ddns-watchdog       // pikaur
```

## 文档资源

> DNSPod API [https://www.dnspod.cn/docs/index.html](https://www.dnspod.cn/docs/index.html)

> Aliyun SDK [GitHub](https://github.com/aliyun/alibaba-cloud-sdk-go) or [https://help.aliyun.com/product/29697.html](https://help.aliyun.com/product/29697.html)

> Cloudflare API [https://api.cloudflare.com/#dns-records-for-a-zone-properties](https://api.cloudflare.com/#dns-records-for-a-zone-properties)
