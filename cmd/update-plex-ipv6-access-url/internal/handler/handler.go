package handler

import (
	"net"
	"net/netip"
	"net/url"
	"strings"

	"github.com/cetteup/update-plex-ipv6-access-url/internal/plex"
)

func UpdateIPv6CustomAccessURL(serverAddr string, token string, interfaceAddr netip.Addr) error {
	client := plex.NewApiClient(serverAddr, token)
	identity, err := client.GetIdentity()
	if err != nil {
		return err
	}

	plexDirectHostname, err := getPlexDirectHostname(token, identity.MachineIdentifier)
	if err != nil {
		return err
	}

	preferences, err := client.GetPreferences()
	if err != nil {
		return err
	}

	currentAccessURLs, err := getCustomConnections(preferences)
	if err != nil {
		return err
	}

	mappedPort, err := getMappedPort(preferences)
	if err != nil {
		return err
	}

	// Drop any existing IPv6 custom access urls (and empty ones) before adding a new one
	targetAccessURLs := make([]string, 0)
	for _, c := range currentAccessURLs {
		drop, err := isIPv6CustomAccessURL(c)
		if err != nil {
			return err
		}

		if !drop && c != "" {
			targetAccessURLs = append(targetAccessURLs, c)
		}
	}

	targetAccessURLs = append(targetAccessURLs, buildIPv6CustomAccessURL(interfaceAddr, plexDirectHostname, mappedPort))

	return client.UpdateCustomConnections(strings.Join(targetAccessURLs, ","))
}

func getPlexDirectHostname(token string, identifier string) (string, error) {
	client := plex.NewApiClient(plex.BaseURL, token)
	resources, err := client.GetResources()
	if err != nil {
		return "", err
	}

	device, err := resources.GetDeviceByIdentifier(identifier)
	if err != nil {
		return "", err
	}

	return device.GetPlexDirectHostname()
}

func getCustomConnections(preferences plex.PreferencesDTO) ([]string, error) {
	setting, err := preferences.GetSettingByID(plex.SettingIDCustomConnections)
	if err != nil {
		return nil, err
	}

	return strings.Split(setting.Value, ","), nil
}

func getMappedPort(preferences plex.PreferencesDTO) (string, error) {
	manualPortMappingMode, err := preferences.GetSettingByID(plex.SettingIDManualPortMappingMode)
	if err != nil {
		return "", err
	}

	var portSetting plex.SettingDTO
	if manualPortMappingMode.IsEnabledBoolSetting() {
		portSetting, err = preferences.GetSettingByID(plex.SettingIDManualPortMappingPort)
		if err != nil {
			return "", err
		}
	} else {
		portSetting, err = preferences.GetSettingByID(plex.SettingIDLastAutomaticMappedPort)
		if err != nil {
			return "", err
		}
	}

	return portSetting.Value, nil
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
