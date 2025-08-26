package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"net"
	"sort"
)

// IRange 接口定义了IP范围的基本操作
type IRange interface {
	ToIp() net.IP           // 如果可以表示为单个IP则返回，否则返回nil
	ToIpNets() []*net.IPNet // 转换为IP网络列表
	ToRange() *Range        // 转换为范围
	String() string         // 字符串表示
}

// Range 表示IP地址范围
type Range struct {
	start net.IP
	end   net.IP
}

// NewRange 创建新的IP范围
func NewRange(start, end net.IP) *Range {
	return &Range{start: start, end: end}
}

// familyLength 返回IP地址族的长度
func (r *Range) familyLength() int {
	return len(r.start)
}

// ToIp 如果范围是单个IP则返回，否则返回nil
func (r *Range) ToIp() net.IP {
	if bytes.Equal(r.start, r.end) {
		return r.start
	}
	return nil
}

// ToIpNets 将IP范围转换为CIDR网络列表
func (r *Range) ToIpNets() []*net.IPNet {
	s, end := r.start, r.end
	ipBits := len(s) * 8

	var result []*net.IPNet
	for {
		if bytes.Compare(s, end) > 0 {
			break
		}

		// 计算最大前缀长度
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

// ToRange 返回自身
func (r *Range) ToRange() *Range {
	return r
}

// String 返回范围的字符串表示
func (r *Range) String() string {
	return ipToString(r.start) + "-" + ipToString(r.end)
}

// IpWrapper 包装单个IP地址
type IpWrapper struct {
	net.IP
}

// ToIp 返回包装的IP
func (r IpWrapper) ToIp() net.IP {
	return r.IP
}

// ToIpNets 将单个IP转换为CIDR网络
func (r IpWrapper) ToIpNets() []*net.IPNet {
	ipBits := len(r.IP) * 8
	return []*net.IPNet{
		{IP: r.IP, Mask: net.CIDRMask(ipBits, ipBits)},
	}
}

// ToRange 将单个IP转换为范围
func (r IpWrapper) ToRange() *Range {
	return &Range{start: r.IP, end: r.IP}
}

// String 返回IP的字符串表示
func (r IpWrapper) String() string {
	return ipToString(r.IP)
}

// IpNetWrapper 包装IP网络
type IpNetWrapper struct {
	*net.IPNet
}

// ToIp 如果网络掩码全为1则返回IP，否则返回nil
func (r IpNetWrapper) ToIp() net.IP {
	if allFF(r.IPNet.Mask) {
		return r.IPNet.IP
	}
	return nil
}

// ToIpNets 返回包装的网络
func (r IpNetWrapper) ToIpNets() []*net.IPNet {
	return []*net.IPNet{r.IPNet}
}

// ToRange 将网络转换为范围
func (r IpNetWrapper) ToRange() *Range {
	ipNet := r.IPNet
	return &Range{start: ipNet.IP, end: lastIp(ipNet)}
}

// String 返回网络的字符串表示
func (r IpNetWrapper) String() string {
	return r.IPNet.String()
}

// 辅助函数
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

// ParseIPRange 解析IP范围字符串
func ParseIPRange(s string) (IRange, error) {
	// 尝试解析为单个IP
	if ip := net.ParseIP(s); ip != nil {
		return IpWrapper{ip}, nil
	}

	// 尝试解析为CIDR
	if _, ipNet, err := net.ParseCIDR(s); err == nil {
		return IpNetWrapper{ipNet}, nil
	}

	// 尝试解析为IP范围 (start-end)
	if idx := bytes.IndexByte([]byte(s), '-'); idx > 0 {
		start := net.ParseIP(s[:idx])
		end := net.ParseIP(s[idx+1:])
		if start != nil && end != nil {
			return NewRange(start, end), nil
		}
	}

	return nil, fmt.Errorf("无法解析IP范围: %s", s)
}

// MergeRanges 合并多个IP范围
func MergeRanges(ranges []IRange) []IRange {
	if len(ranges) == 0 {
		return ranges
	}

	// 转换为范围并排序
	var rangeList []*Range
	for _, r := range ranges {
		rangeList = append(rangeList, r.ToRange())
	}

	sort.Slice(rangeList, func(i, j int) bool {
		return bytes.Compare(rangeList[i].start, rangeList[j].start) < 0
	})

	// 合并重叠的范围
	var result []IRange
	current := rangeList[0]

	for i := 1; i < len(rangeList); i++ {
		next := rangeList[i]

		// 检查是否重叠或相邻
		if bytes.Compare(addOne(current.end), next.start) >= 0 {
			// 合并范围
			if bytes.Compare(current.end, next.end) < 0 {
				current.end = next.end
			}
		} else {
			// 添加当前范围并移动到下一个
			result = append(result, current)
			current = next
		}
	}

	result = append(result, current)
	return result
}

// CalculateIPv6CIDRRange 计算两个IPv6地址之间的CIDR范围
func CalculateIPv6CIDRRange(startIP, endIP string) ([]string, error) {
	start := net.ParseIP(startIP)
	end := net.ParseIP(endIP)

	if start == nil || end == nil {
		return nil, fmt.Errorf("无效的IP地址")
	}

	ipRange := NewRange(start, end)
	cidrs := ipRange.ToIpNets()

	var result []string
	for _, cidr := range cidrs {
		result = append(result, cidr.String())
	}

	return result, nil
}
