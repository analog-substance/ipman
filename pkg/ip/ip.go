package ip

import (
	"bytes"
	"net"
	"sort"

	"github.com/analog-substance/ipman/internal/set"
	"github.com/apparentlymart/go-cidr/cidr"
)

type IPSet struct {
	set set.Set
}

func NewSet() IPSet {
	return IPSet{
		set: set.NewSet(""),
	}
}

func (s *IPSet) Add(ip net.IP) bool {
	return s.set.Add(ip.String())
}

func (s *IPSet) AddNetwork(network *net.IPNet) int {
	return s.AddRange(GetIPs(network))
}
func (s *IPSet) AddNetworkWithFilter(network *net.IPNet, predicate func(net.IP) bool) int {
	return s.AddRange(GetIPsWithFilter(network, predicate))
}

func (s *IPSet) AddRange(items []net.IP) int {
	num := 0
	for _, ip := range items {
		if s.Add(ip) {
			num++
		}
	}
	return num
}

func (s *IPSet) Slice() []net.IP {
	strSlice := s.set.Slice().([]string)
	if strSlice == nil {
		return nil
	}

	var ips []net.IP
	for _, v := range strSlice {
		ips = append(ips, net.ParseIP(v))
	}
	return ips
}

func (s *IPSet) SortedSlice() []net.IP {
	ips := s.Slice()
	if ips == nil {
		return nil
	}

	sort.Slice(ips, func(i, j int) bool {
		return bytes.Compare(ips[i], ips[j]) < 0
	})
	return ips
}

func GetIPs(network *net.IPNet) []net.IP {
	return GetIPsWithFilter(network, func(i net.IP) bool { return true })
}

func GetIPsWithFilter(network *net.IPNet, predicate func(net.IP) bool) []net.IP {
	var ips []net.IP

	count := cidr.AddressCount(network)

	var ip net.IP
	for i := uint64(0); i < count; i++ {
		if ip == nil {
			ip = network.IP
		} else {
			ip = cidr.Inc(ip)
		}

		if predicate(ip) {
			ips = append(ips, ip)
		}
	}

	return ips
}
