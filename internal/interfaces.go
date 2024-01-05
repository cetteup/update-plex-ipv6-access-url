package internal

import (
	"net"
	"net/netip"
)

func GetGlobalUnicastIPv6AddrsByInterfaceName(name string) ([]netip.Addr, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	return getInterfaceIPv6GlobalUnicastAddrsByInterface(iface)
}

func getInterfaceIPv6GlobalUnicastAddrsByInterface(iface *net.Interface) ([]netip.Addr, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	ipv6GlobalUnicastAddrs := make([]netip.Addr, 0, len(addrs))
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
			ipv6GlobalUnicastAddrs = append(ipv6GlobalUnicastAddrs, addrFromIP)
		}
	}

	return ipv6GlobalUnicastAddrs, nil
}
