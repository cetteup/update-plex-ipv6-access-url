package handler

import (
	"net"
	"net/netip"
	"net/url"
	"strings"

	"github.com/cetteup/update-plex-ipv6-access-url/internal/plex"
)

func UpdateIPv6CustomAccessURL(addr netip.Addr, config plex.Config) error {
	client := plex.NewApiClient(plex.BaseURL, config.Preferences.GetToken())
	resources, err := client.GetResources()
	if err != nil {
		return err
	}

	device, err := resources.GetDeviceByIdentifier(config.Preferences.GetProcessedMachineIdentifier())

	plexDirectHostname, err := device.GetPlexDirectHostname()
	if err != nil {
		return err
	}

	localConnection, err := device.GetLocalConnection()
	if err != nil {
		return err
	}

	// Drop any existing IPv6 custom access urls (and empty ones) before adding a new one
	customAccessURLs := make([]string, 0)
	for _, c := range config.Preferences.GetCustomConnections() {
		drop, err := isIPv6CustomAccessURL(c)
		if err != nil {
			return err
		}

		if !drop && c != "" {
			customAccessURLs = append(customAccessURLs, c)
		}
	}

	customAccessURLs = append(customAccessURLs, buildIPv6CustomAccessURL(addr, plexDirectHostname, config.Preferences.GetMappedPort()))

	localClient := plex.NewApiClient(localConnection.URI, config.Preferences.GetToken())
	return localClient.UpdateCustomConnections(strings.Join(customAccessURLs, ","))
}

func buildIPv6CustomAccessURL(addr netip.Addr, plexDirectHostname, port string) string {
	dashedIPv6 := strings.ReplaceAll(addr.StringExpanded(), ":", "-")
	hostname := strings.Join([]string{dashedIPv6, plexDirectHostname}, ".")

	u := url.URL{
		Host:   net.JoinHostPort(hostname, port),
		Scheme: "https",
	}
	return u.String()
}

func isIPv6CustomAccessURL(customAccessURL string) (bool, error) {
	u, err := url.Parse(customAccessURL)
	if err != nil {
		return false, err
	}

	hostname := u.Hostname()
	if !strings.Contains(hostname, ".plex.direct") {
		return false, nil
	}

	hostElems := strings.Split(hostname, ".")
	if len(hostElems) != 4 {
		return false, nil
	}

	ipElems := strings.Split(hostElems[0], "-")
	if len(ipElems) != 8 {
		return false, err
	}

	addr, err := netip.ParseAddr(strings.Join(ipElems, ":"))
	if err != nil {
		return false, err
	}

	return addr.Is6(), nil
}
