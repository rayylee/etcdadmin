package utils

import (
	"net"
)

func IsLegalIP4(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			address := net.ParseIP(s)
			if address != nil {
				return true
			}
		}
	}

	return false
}

func GetHostIP4() ([]string, error) {
	ifaces, err := net.Interfaces()
	var ips []string
	if err != nil {
		return ips, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := getIp4FromAddr(addr)
			if ip == nil {
				continue
			}
			ips = append(ips, ip.String())
		}
	}
	return ips, err
}

func getIp4FromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
