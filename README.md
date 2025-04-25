# realIP

[English](README_EN.md) | 简体中文

一个高效的命令行工具，用于筛选目标真实使用IP。

## 核心逻辑
- 过滤域名直接解析到多个IP（默认支持）
- 过滤域名解析到CDN CNAME（参数支持）
- 最终目标只获取直接A记录解析到单个IP

## 功能特点

- 支持单个域名解析或批量解析域名列表
- 支持自定义输出格式（JSON、CSV）
- 支持多线程并发处理
- 提供完整的日志输出
- 支持自定义输出位置
- 支持过滤CNAME记录，只获取直接解析到IP的域名

## 安装

### 从源码编译

```bash
git clone https://github.com/r00tuser111/realIP.git
cd realIP
go build -o realip
```

## 使用方法

### 命令行参数

- `-target`, `-u`: 指定单个目标域名
- `-list`, `-l`: 指定包含多个域名的文件路径（每行一个域名）
- `-json`: 以JSON格式输出结果
- `-csv`: 以CSV格式输出结果
- `-threads`, `-t`: 设置并发线程数（默认为10）
- `-o`: 指定输出文件路径
- `-silent`: 关闭所有日志输出
- `-cname`: 过滤CNAME记录，只显示直接解析到IP的域名

### 示例

解析单个域名:
```bash
./realip -u example.com
```

解析域名列表:
```bash
./realip -l domains.txt
```

以JSON格式输出并保存到文件:
```bash
./realip -u example.com -json -o results.json
```

以CSV格式输出并保存到文件:
```bash
./realip -l domains.txt -csv -o results.csv -t 20
```

静默模式运行:
```bash
./realip -l domains.txt -o results.txt -silent
```

过滤CNAME记录，只获取直接解析到IP的域名:
```bash
./realip -l domains.txt -cname -o results.txt
```

### 输出格式

#### 文本格式（默认）
```
域名           IP地址           CNAME    状态
----------------------------------------
example.com    93.184.216.34   否       成功
cdn.site.com   -               是       失败
```

#### JSON格式
```json
[
  {
    "domain": "example.com",
    "ip": "93.184.216.34",
    "has_cname": false,
    "timestamp": "2024-03-18T15:04:05Z",
    "error": ""
  }
]
```

#### CSV格式
```csv
Domain,IP,Has_CNAME,Timestamp,Error
example.com,93.184.216.34,false,2024-03-18T15:04:05Z,
cdn.site.com,,true,2024-03-18T15:04:06Z,"域名存在CNAME记录"
```

## 许可证

MIT