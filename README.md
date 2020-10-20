## dnsctl

> dnsctl 是一个简单的命令行工具用于控制 CoreDNS 的 [`etcdhost`](https://github.com/ytpay/etcdhosts) 插件，从而实现对 etcd 内 hosts 文件一处修改全局生效。

### 一、配置

默认情况下 dnsctl 读取 `$HOME/.dnsctl.yaml` 配置来链接 etcd 集群，配置样例如下:

```yaml
# etcd 中 etcdhosts 插件的 key
dnskey: /etcdhosts
# etcd 集群配置
etcd:
  cert: /etc/etcd/ssl/etcd.pem
  key: /etc/etcd/ssl/etcd-key.pem
  ca: /etc/etcd/ssl/etcd-root-ca.pem
  endpoints:
    - https://172.16.10.11:2379
    - https://172.16.10.12:2379
    - https://172.16.10.13:2379
```

**可通过 `dnsctl config` 命令查看样例配置**

### 二、使用

```sh
➜ ~ dnsctl --help
dnsctl for etcdhosts plugin

Usage:
  dnsctl [flags]
  dnsctl [command]

Available Commands:
  config      show example config
  dump        dump hosts
  edit        edit hosts
  help        Help about any command
  upload      upload hosts from file
  version     show hosts version

Flags:
      --config string   config file (default is $HOME/.dnsctl.yaml)
  -h, --help            help for dnsctl
      --version         version for dnsctl

Use "dnsctl [command] --help" for more information about a command.
```
