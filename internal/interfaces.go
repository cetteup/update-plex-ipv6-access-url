package internal

import (
	"fmt"
	"net"
	"net/netip"
)

func GetInterfaceGlobalUnicastIPv6ByName(name string) (netip.Addr, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return netip.Addr{}, err
	}

	return getInterfaceIPv6GlobalUnicastIP(iface)
}

func getInterfaceIPv6GlobalUnicastIP(iface *net.Interface) (netip.Addr, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return netip.Addr{}, err
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPAddr:
			ip = v.IP
		case *net.IPNet:
			ip = v.IP
		default:
			continue
		}

		if addrFromIP, ok := netip.AddrFromSlice(ip); ok && addrFromIP.Is6() && !addrFromIP.Is4In6() && addrFromIP.IsGlobalUnicast() && !addrFromIP.IsPrivate() {
			return addrFromIP, nil
		}
	}

	return netip.Addr{}, fmt.Errorf("no global unicast IPv6 address found on interface: %s", iface.Name)
}
