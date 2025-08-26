package main

import (
	"fmt"
	"net"

	gocidrpkg "github.com/liliangshan/go-cidr-pkg"
)

func main() {
	startIP := "2001:200:141::"
	endIP := "2001:200:142:ffff:ffff:ffff:ffff:ffff"

	fmt.Printf("Using cidrmerger library to calculate IPv6 address range:\n")
	fmt.Printf("Start address: %s\n", startIP)
	fmt.Printf("End address: %s\n", endIP)
	fmt.Println()

	// Method 1: Using CalculateIPv6CIDRRange function
	cidrs, err := gocidrpkg.CalculateIPv6CIDRRange(startIP, endIP)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("CIDR ranges:\n")
	for i, cidr := range cidrs {
		fmt.Printf("%d. %s\n", i+1, cidr)
	}

	fmt.Printf("\nTotal generated %d CIDR blocks\n", len(cidrs))

	// Method 2: Using Range object
	fmt.Println("\nUsing Range object:")
	start := net.ParseIP(startIP)
	end := net.ParseIP(endIP)

	ipRange := gocidrpkg.NewRange(start, end)
	fmt.Printf("IP range: %s\n", ipRange.String())

	cidrNets := ipRange.ToIpNets()
	fmt.Printf("CIDR network count: %d\n", len(cidrNets))

	// Method 3: Parse IP range string
	fmt.Println("\nParsing IP range string:")
	rangeStr := startIP + "-" + endIP
	if parsedRange, err := gocidrpkg.ParseIPRange(rangeStr); err == nil {
		fmt.Printf("Parse successful: %s\n", parsedRange.String())
		cidrs2 := parsedRange.ToIpNets()
		fmt.Printf("Generated CIDR count: %d\n", len(cidrs2))
	} else {
		fmt.Printf("Parse failed: %v\n", err)
	}

	// Test IPv4
	fmt.Println("\nTesting IPv4 address range:")
	ipv4Start := "192.168.1.0"
	ipv4End := "192.168.1.255"

	ipv4Range := gocidrpkg.NewRange(net.ParseIP(ipv4Start), net.ParseIP(ipv4End))
	ipv4Cidrs := ipv4Range.ToIpNets()

	fmt.Printf("IPv4 range %s-%s converted to %d CIDR blocks:\n", ipv4Start, ipv4End, len(ipv4Cidrs))
	for i, cidr := range ipv4Cidrs {
		fmt.Printf("%d. %s\n", i+1, cidr.String())
	}
}
