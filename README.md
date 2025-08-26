# go-cidr-pkg

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A powerful Go library for handling IP address ranges, CIDR networks, and IP address merging.

## ✨ Features

- 🚀 **IPv4/IPv6 Dual Stack Support** - Full support for IPv4 and IPv6 address formats
- 🔧 **Smart CIDR Calculation** - Convert IP address ranges to optimal CIDR blocks
- 📝 **Multi-format Parsing** - Support for single IP, CIDR, IP range and other input formats
- 🔄 **Range Merging** - Intelligently merge overlapping IP ranges
- 💻 **Pure Go Implementation** - No external dependencies, excellent performance

## 🚀 Quick Start

### Installation

```bash
go get github.com/liliangshan/go-cidr-pkg
```

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/liliangshan/go-cidr-pkg"
)

func main() {
    // Calculate CIDR for IPv6 address range
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

## 📚 Usage Examples

### IPv6 CIDR Calculation

```go
import "github.com/liliangshan/go-cidr-pkg"

// Calculate IPv6 address range
cidrs, err := gocidrpkg.CalculateIPv6CIDRRange(
    "2001:200:141::", 
    "2001:200:142:ffff:ffff:ffff:ffff:ffff"
)

// Output: ["2001:200:141::/48", "2001:200:142::/47"]
```

### Using Range Object

```go
import (
    "net"
    "github.com/liliangshan/go-cidr-pkg"
)

start := net.ParseIP("192.168.1.0")
end := net.ParseIP("192.168.1.255")

ipRange := gocidrpkg.NewRange(start, end)
cidrNets := ipRange.ToIpNets()

// Output: [192.168.1.0/24]
```

## 🏗️ Project Structure

```
go-cidr-pkg/
├── go-cidr-pkg.go        # Core library file
├── example/               # Usage examples
│   └── main.go           # Example code
├── README.md              # Project documentation
├── LICENSE                # MIT License
├── go.mod                 # Go module configuration
└── .gitignore             # Git ignore file
```

## 🔧 API Reference

### Main Functions

| Function | Description |
|----------|-------------|
| `NewRange(start, end net.IP) *Range` | Create new IP range |
| `ParseIPRange(s string) (IRange, error)` | Parse IP range string |
| `MergeRanges(ranges []IRange) []IRange` | Merge multiple IP ranges |
| `CalculateIPv6CIDRRange(startIP, endIP string) ([]string, error)` | Calculate IPv6 CIDR range |

### Types

- `IRange` - IP range interface
- `Range` - IP address range structure
- `IpWrapper` - Single IP wrapper
- `IpNetWrapper` - IP network wrapper

## 🧪 Testing

Run examples:

```bash
go run example/main.go
```

## 📄 License

This project is open source under the MIT License. See the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Welcome to submit Issues and Pull Requests!

---

If this project helps you, please give us a ⭐ Star!
