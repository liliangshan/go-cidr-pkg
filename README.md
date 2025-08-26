# go-cidr-pkg

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

一个强大的Go语言库，用于处理IP地址范围、CIDR网络和IP地址合并。基于原始 `github.com/zhanhb/cidr-merger` 项目重构，提供更友好的Go API接口。

## ✨ 功能特性

- 🚀 **IPv4/IPv6 双栈支持** - 完全支持IPv4和IPv6地址格式
- 🔧 **智能CIDR计算** - 将IP地址范围转换为最优CIDR块
- 📝 **多格式解析** - 支持单个IP、CIDR、IP范围等多种输入格式
- 🔄 **范围合并** - 智能合并重叠的IP范围
- 💻 **纯Go实现** - 无外部依赖，性能优秀

## 🚀 快速开始

### 安装

```bash
go get github.com/liliangshan/go-cidr-pkg
```

### 基本使用

```go
package main

import (
    "fmt"
    "github.com/liliangshan/go-cidr-pkg"
)

func main() {
    // 计算IPv6地址范围的CIDR
    startIP := "2001:200:141::"
    endIP := "2001:200:142:ffff:ffff:ffff:ffff:ffff"
    
    cidrs, err := gocidrpkg.CalculateIPv6CIDRRange(startIP, endIP)
    if err != nil {
        panic(err)
    }
    
    for i, cidr := range cidrs {
        fmt.Printf("%d. %s\n", i+1, cidr)
    }
}
```

## 📚 使用示例

### IPv6 CIDR 计算

```go
import "github.com/liliangshan/go-cidr-pkg"

// 计算IPv6地址范围
cidrs, err := gocidrpkg.CalculateIPv6CIDRRange(
    "2001:200:141::", 
    "2001:200:142:ffff:ffff:ffff:ffff:ffff"
)

// 输出: ["2001:200:141::/48", "2001:200:142::/47"]
```

### 使用Range对象

```go
import (
    "net"
    "github.com/liliangshan/go-cidr-pkg"
)

start := net.ParseIP("192.168.1.0")
end := net.ParseIP("192.168.1.255")

ipRange := gocidrpkg.NewRange(start, end)
cidrNets := ipRange.ToIpNets()

// 输出: [192.168.1.0/24]
```

## 🏗️ 项目结构

```
go-cidr-pkg/
├── go-cidr-pkg.go        # 核心库文件
├── example/               # 使用示例
│   └── main.go           # 示例代码
├── README.md              # 项目文档
├── LICENSE                # MIT许可证
├── go.mod                 # Go模块配置
└── .gitignore             # Git忽略文件
```

## 🔧 API 参考

### 主要函数

| 函数 | 描述 |
|------|------|
| `NewRange(start, end net.IP) *Range` | 创建新的IP范围 |
| `ParseIPRange(s string) (IRange, error)` | 解析IP范围字符串 |
| `MergeRanges(ranges []IRange) []IRange` | 合并多个IP范围 |
| `CalculateIPv6CIDRRange(startIP, endIP string) ([]string, error)` | 计算IPv6 CIDR范围 |

### 类型

- `IRange` - IP范围接口
- `Range` - IP地址范围结构
- `IpWrapper` - 单个IP包装器
- `IpNetWrapper` - IP网络包装器

## 🧪 测试

运行示例：

```bash
go run example/main.go
```

## 📄 许可证

本项目基于 MIT 许可证开源。详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

如果这个项目对您有帮助，请给我们一个 ⭐ Star！
