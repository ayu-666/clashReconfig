# clash多订阅整合根据

[官方文档-配置文件示例](https://github.com/Dreamacro/clash/wiki/configuration#all-configuration-options)

### 配置

将订阅地址配置到clash配置文件proxy-providers中，保存为`模板.yaml`

```yaml
proxy-groups:
  - name: "LoadBalance"
    type: "load-balance"
    url: 'http://www.gstatic.com/generate_204'
    interval: 300
    use:
      - provider01
proxy-providers:
  provider01:
    type: http
    url: "https://ghproxy.com/raw.githubusercontent.com/ayu-666/clashReconfig/main/examples/proxy_provider1.yaml"
    interval: 3600
    path: pp1.yaml
```

### 命令


```bash
clashReconfig -i=模板.yaml -o=输出.yaml -s=3600
```

- `-i` 输入文件路径
- `-o` 输出文件路径
- `-k` 保持进程，根据配置文件的interval属性定时更新订阅，同时监听`模板.yaml`变化
- `-c` 每次输出文件后执行命令

### 输出

```yaml
proxies:
  - name: "provider1::px1"
    type: ss
    server: 10.0.0.1
    port: 80
    cipher: aes-256-cfb
    password: "password123456"
    interval: 300
proxy-groups:
  - name: "LoadBalance"
    type: "load-balance"
    url: 'http://www.gstatic.com/generate_204'
    interval: 300
    proxies:
      - provider1::px1
```
