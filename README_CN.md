# Go CIDR 包

一个用于处理 IP 地址范围、CIDR 网络和 IP 地址操作的 Go 语言包。

## 功能特性

- 🚀 **IP 范围处理**：支持 IP 地址范围的解析和操作
- 🌐 **CIDR 网络支持**：完整的 CIDR 网络地址处理
- 🔄 **范围合并**：自动合并重叠的 IP 范围
- 📱 **IPv4/IPv6 支持**：同时支持 IPv4 和 IPv6 地址
- 🎯 **灵活解析**：支持多种 IP 地址格式的解析

## 安装

```bash
go get github.com/liliangshan/go-cidr-pkg
```

## 快速开始

### 基本用法

```go
package main

import (
    "fmt"
    "log"
    "github.com/liliangshan/go-cidr-pkg"
)

func main() {
    // 解析单个 IP 地址
    ipRange, err := gocidrpkg.ParseIPRange("192.168.1.1")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("IP: %s\n", ipRange.String())

    // 解析 CIDR 网络
    cidrRange, err := gocidrpkg.ParseIPRange("192.168.1.0/24")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("CIDR: %s\n", cidrRange.String())

    // 解析 IP 范围
    rangeRange, err := gocidrpkg.ParseIPRange("192.168.1.1-192.168.1.100")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("范围: %s\n", rangeRange.String())
}
```

## API 参考

### 核心接口

#### IRange 接口

```go
type IRange interface {
    ToIp() net.IP                    // 转换为单个 IP（如果可能）
    ToIpNets() []*net.IPNet         // 转换为 CIDR 网络列表
    ToRange() *Range                // 转换为 IP 范围
    String() string                 // 字符串表示
}
```

### 主要类型

#### Range 结构体

表示一个 IP 地址范围：

```go
type Range struct {
    start net.IP  // 起始 IP
    end   net.IP  // 结束 IP
}

// 创建新的 IP 范围
func NewRange(start, end net.IP) *Range
```

#### IpWrapper 结构体

包装单个 IP 地址：

```go
type IpWrapper struct {
    net.IP
}
```

#### IpNetWrapper 结构体

包装 CIDR 网络：

```go
type IpNetWrapper struct {
    *net.IPNet
}
```

### 主要函数

#### ParseIPRange

解析 IP 范围字符串：

```go
func ParseIPRange(s string) (IRange, error)
```

支持的格式：
- 单个 IP：`"192.168.1.1"`
- CIDR 网络：`"192.168.1.0/24"`
- IP 范围：`"192.168.1.1-192.168.1.100"`

#### MergeRanges

合并多个 IP 范围：

```go
func MergeRanges(ranges []IRange) []IRange
```

#### CalculateIPv6CIDRRange

计算两个 IPv6 地址之间的 CIDR 范围：

```go
func CalculateIPv6CIDRRange(startIP, endIP string) ([]string, error)
```

## 使用示例

### 示例 1：解析不同类型的 IP 地址

```go
func example1() {
    // 解析单个 IP
    ip, _ := gocidrpkg.ParseIPRange("10.0.0.1")
    fmt.Println("单个 IP:", ip.String())

    // 解析 CIDR
    cidr, _ := gocidrpkg.ParseIPRange("10.0.0.0/16")
    fmt.Println("CIDR:", cidr.String())

    // 解析范围
    rng, _ := gocidrpkg.ParseIPRange("10.0.0.1-10.0.0.10")
    fmt.Println("范围:", rng.String())
}
```

### 示例 2：合并 IP 范围

```go
func example2() {
    ranges := []gocidrpkg.IRange{}
    
    // 添加一些重叠的范围
    r1, _ := gocidrpkg.ParseIPRange("192.168.1.1-192.168.1.50")
    r2, _ := gocidrpkg.ParseIPRange("192.168.1.30-192.168.1.100")
    r3, _ := gocidrpkg.ParseIPRange("192.168.2.1-192.168.2.50")
    
    ranges = append(ranges, r1, r2, r3)
    
    // 合并范围
    merged := gocidrpkg.MergeRanges(ranges)
    
    for _, r := range merged {
        fmt.Println("合并后:", r.String())
    }
}
```

### 示例 3：IPv6 地址处理

```go
func example3() {
    startIP := "2001:db8::1"
    endIP := "2001:db8::100"
    
    cidrs, err := gocidrpkg.CalculateIPv6CIDRRange(startIP, endIP)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("IPv6 CIDR 范围:")
    for _, cidr := range cidrs {
        fmt.Println("  ", cidr)
    }
}
```

## 错误处理

包会返回详细的错误信息：

```go
ipRange, err := gocidrpkg.ParseIPRange("invalid-ip")
if err != nil {
    fmt.Printf("解析错误: %v\n", err)
    return
}
```

常见错误：
- `unable to parse IP range: invalid-ip` - 无法解析的 IP 地址格式
- `invalid IP address` - 无效的 IP 地址

## 性能考虑

- 所有操作都在内存中完成，适合中小规模的 IP 地址处理
- 对于大量 IP 地址，建议分批处理
- IPv6 地址处理可能需要更多内存

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 更新日志

### v1.0.0
- 初始版本发布
- 支持 IPv4 和 IPv6 地址
- 完整的 IP 范围操作功能
- CIDR 网络支持
