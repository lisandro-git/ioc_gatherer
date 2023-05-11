package cmd

import (
	"net"
	"encoding/binary"
	"errors"
)

var (
	ErrNotIPv4Address = errors.New("not an IPv4 addres")
)

func IPv4ToInt(ipaddr net.IP) (uint32, error) {
	if ipaddr.To4() == nil {
		return 0, ErrNotIPv4Address
	}
	return binary.BigEndian.Uint32(ipaddr.To4()), nil
}

func IntToIPv4(ipaddr uint32) net.IP {
	ip := make(net.IP, net.IPv4len)

	// Proceed conversion
	binary.BigEndian.PutUint32(ip, ipaddr)

	return ip
}
