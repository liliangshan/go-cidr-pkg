package main

import (
	"fmt"
	"net"

	gocidrpkg "github.com/liliangshan/go-cidr-pkg"
)

func main() {
	startIP := "2001:200:141::"
	endIP := "2001:200:142:ffff:ffff:ffff:ffff:ffff"

	fmt.Printf("使用cidrmerger库计算IPv6地址范围:\n")
	fmt.Printf("起始地址: %s\n", startIP)
	fmt.Printf("结束地址: %s\n", endIP)
	fmt.Println()

	// 方法1: 使用CalculateIPv6CIDRRange函数
	cidrs, err := gocidrpkg.CalculateIPv6CIDRRange(startIP, endIP)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	fmt.Printf("CIDR范围:\n")
	for i, cidr := range cidrs {
		fmt.Printf("%d. %s\n", i+1, cidr)
	}

	fmt.Printf("\n总共生成了 %d 个CIDR块\n", len(cidrs))

	// 方法2: 使用Range对象
	fmt.Println("\n使用Range对象:")
	start := net.ParseIP(startIP)
	end := net.ParseIP(endIP)

	ipRange := gocidrpkg.NewRange(start, end)
	fmt.Printf("IP范围: %s\n", ipRange.String())

	cidrNets := ipRange.ToIpNets()
	fmt.Printf("CIDR网络数量: %d\n", len(cidrNets))

	// 方法3: 解析IP范围字符串
	fmt.Println("\n解析IP范围字符串:")
	rangeStr := startIP + "-" + endIP
	if parsedRange, err := gocidrpkg.ParseIPRange(rangeStr); err == nil {
		fmt.Printf("解析成功: %s\n", parsedRange.String())
		cidrs2 := parsedRange.ToIpNets()
		fmt.Printf("生成的CIDR数量: %d\n", len(cidrs2))
	} else {
		fmt.Printf("解析失败: %v\n", err)
	}

	// 测试IPv4
	fmt.Println("\n测试IPv4地址范围:")
	ipv4Start := "192.168.1.0"
	ipv4End := "192.168.1.255"

	ipv4Range := gocidrpkg.NewRange(net.ParseIP(ipv4Start), net.ParseIP(ipv4End))
	ipv4Cidrs := ipv4Range.ToIpNets()

	fmt.Printf("IPv4范围 %s-%s 转换为 %d 个CIDR块:\n", ipv4Start, ipv4End, len(ipv4Cidrs))
	for i, cidr := range ipv4Cidrs {
		fmt.Printf("%d. %s\n", i+1, cidr.String())
	}
}
