package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"net"
	"sort"
)

// IRange interface defines basic operations for IP ranges
type IRange interface {
	ToIp() net.IP           // Returns IP if it can be represented as a single IP, otherwise returns nil
	ToIpNets() []*net.IPNet // Converts to IP network list
	ToRange() *Range        // Converts to range
	String() string         // String representation
}

// Range represents an IP address range
type Range struct {
	start net.IP
	end   net.IP
}

// NewRange creates a new IP range
func NewRange(start, end net.IP) *Range {
	return &Range{start: start, end: end}
}

// familyLength returns the length of the IP address family
func (r *Range) familyLength() int {
	return len(r.start)
}

// ToIp returns the IP if the range is a single IP, otherwise returns nil
func (r *Range) ToIp() net.IP {
	if bytes.Equal(r.start, r.end) {
		return r.start
	}
	return nil
}

// ToIpNets converts IP range to CIDR network list
func (r *Range) ToIpNets() []*net.IPNet {
	s, end := r.start, r.end
	ipBits := len(s) * 8

	var result []*net.IPNet
	for {
		if bytes.Compare(s, end) > 0 {
			break
		}

		// Calculate maximum prefix length
		cidr := max(prefixLength(xor(addOne(end), s)), ipBits-trailingZeros(s))
		ipNet := &net.IPNet{IP: s, Mask: net.CIDRMask(cidr, ipBits)}
		result = append(result, ipNet)

		tmp := lastIp(ipNet)
		if !lessThan(tmp, end) {
			break
		}
		s = addOne(tmp)
	}
	return result
}

// ToRange returns itself
func (r *Range) ToRange() *Range {
	return r
}

// String returns the string representation of the range
func (r *Range) String() string {
	return ipToString(r.start) + "-" + ipToString(r.end)
}

// IpWrapper wraps a single IP address
type IpWrapper struct {
	net.IP
}

// ToIp returns the wrapped IP
func (r IpWrapper) ToIp() net.IP {
	return r.IP
}

// ToIpNets converts a single IP to CIDR network
func (r IpWrapper) ToIpNets() []*net.IPNet {
	ipBits := len(r.IP) * 8
	return []*net.IPNet{
		{IP: r.IP, Mask: net.CIDRMask(ipBits, ipBits)},
	}
}

// ToRange converts a single IP to range
func (r IpWrapper) ToRange() *Range {
	return &Range{start: r.IP, end: r.IP}
}

// String returns the string representation of the IP
func (r IpWrapper) String() string {
	return ipToString(r.IP)
}

// IpNetWrapper wraps IP network
type IpNetWrapper struct {
	*net.IPNet
}

// ToIp returns IP if network mask is all 1s, otherwise returns nil
func (r IpNetWrapper) ToIp() net.IP {
	if allFF(r.IPNet.Mask) {
		return r.IPNet.IP
	}
	return nil
}

// ToIpNets returns the wrapped network
func (r IpNetWrapper) ToIpNets() []*net.IPNet {
	return []*net.IPNet{r.IPNet}
}

// ToRange converts network to range
func (r IpNetWrapper) ToRange() *Range {
	ipNet := r.IPNet
	return &Range{start: ipNet.IP, end: lastIp(ipNet)}
}

// String returns the string representation of the network
func (r IpNetWrapper) String() string {
	return r.IPNet.String()
}

// Helper functions
func ipToString(ip net.IP) string {
	if len(ip) == net.IPv6len {
		if ipv4 := ip.To4(); len(ipv4) == net.IPv4len {
			return "::ffff:" + ipv4.String()
		}
	}
	return ip.String()
}

func addOne(ip net.IP) net.IP {
	result := make(net.IP, len(ip))
	copy(result, ip)

	for i := len(result) - 1; i >= 0; i-- {
		result[i]++
		if result[i] != 0 {
			break
		}
	}
	return result
}

func lastIp(ipNet *net.IPNet) net.IP {
	ip := make(net.IP, len(ipNet.IP))
	copy(ip, ipNet.IP)

	mask := ipNet.Mask
	for i := 0; i < len(ip); i++ {
		ip[i] |= ^mask[i]
	}
	return ip
}

func prefixLength(x net.IP) int {
	for i := 0; i < len(x); i++ {
		if x[i] != 0 {
			return i*8 + bits.LeadingZeros8(x[i])
		}
	}
	return len(x) * 8
}

func trailingZeros(ip net.IP) int {
	for i := len(ip) - 1; i >= 0; i-- {
		if ip[i] != 0 {
			return (len(ip)-1-i)*8 + bits.TrailingZeros8(ip[i])
		}
	}
	return len(ip) * 8
}

func xor(a, b net.IP) net.IP {
	result := make(net.IP, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}
	return result
}

func lessThan(a, b net.IP) bool {
	return bytes.Compare(a, b) < 0
}

func allFF(mask net.IPMask) bool {
	for _, b := range mask {
		if b != 0xff {
			return false
		}
	}
	return true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ParseIPRange parses IP range string
func ParseIPRange(s string) (IRange, error) {
	// Try to parse as single IP
	if ip := net.ParseIP(s); ip != nil {
		return IpWrapper{ip}, nil
	}

	// Try to parse as CIDR
	if _, ipNet, err := net.ParseCIDR(s); err == nil {
		return IpNetWrapper{ipNet}, nil
	}

	// Try to parse as IP range (start-end)
	if idx := bytes.IndexByte([]byte(s), '-'); idx > 0 {
		start := net.ParseIP(s[:idx])
		end := net.ParseIP(s[idx+1:])
		if start != nil && end != nil {
			return NewRange(start, end), nil
		}
	}

	return nil, fmt.Errorf("unable to parse IP range: %s", s)
}

// MergeRanges merges multiple IP ranges
func MergeRanges(ranges []IRange) []IRange {
	if len(ranges) == 0 {
		return ranges
	}

	// Convert to ranges and sort
	var rangeList []*Range
	for _, r := range ranges {
		rangeList = append(rangeList, r.ToRange())
	}

	sort.Slice(rangeList, func(i, j int) bool {
		return bytes.Compare(rangeList[i].start, rangeList[j].start) < 0
	})

	// Merge overlapping ranges
	var result []IRange
	current := rangeList[0]

	for i := 1; i < len(rangeList); i++ {
		next := rangeList[i]

		// Check if overlapping or adjacent
		if bytes.Compare(addOne(current.end), next.start) >= 0 {
			// Merge ranges
			if bytes.Compare(current.end, next.end) < 0 {
				current.end = next.end
			}
		} else {
			// Add current range and move to next
			result = append(result, current)
			current = next
		}
	}

	result = append(result, current)
	return result
}

// CalculateIPv6CIDRRange calculates CIDR range between two IPv6 addresses
func CalculateIPv6CIDRRange(startIP, endIP string) ([]string, error) {
	start := net.ParseIP(startIP)
	end := net.ParseIP(endIP)

	if start == nil || end == nil {
		return nil, fmt.Errorf("invalid IP address")
	}

	ipRange := NewRange(start, end)
	cidrs := ipRange.ToIpNets()

	var result []string
	for _, cidr := range cidrs {
		result = append(result, cidr.String())
	}

	return result, nil
}
