package plex

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiClient_GetIdentity(t *testing.T) {
	token := "some-token"
	timeout := 5

	tests := []struct {
		name              string
		givenStatusCode   int
		givenData         []byte
		wantIdentity      IdentityDTO
		wantErrorContains string
	}{
		{
			name:            "successfully fetches identity",
			givenStatusCode: 200,
			givenData: []byte(`
				<?xml version="1.0" encoding="UTF-8"?>\n
				<MediaContainer size="0" claimed="1" machineIdentifier="1142ed040a27acc36ea876e8362b28464c3d240d" version="1.31.1.6733-bc0674160"></MediaContainer>
			`),
			wantIdentity: IdentityDTO{
				MachineIdentifier: "1142ed040a27acc36ea876e8362b28464c3d240d",
			},
		},
		{
			name:              "returns error for non-200 response code",
			givenStatusCode:   401,
			wantErrorContains: "failed with status code 401 (401 Unauthorized)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, token, r.Header.Get(headerKeyToken))
				assert.Equal(t, identityEndpoint, r.URL.Path)

				w.WriteHeader(tt.givenStatusCode)
				_, err := w.Write(tt.givenData)
				require.NoError(t, err)
			}))

			client := NewApiClient(server.URL, token, timeout)

			// WHEN
			resources, err := client.GetIdentity()

			// THEN
			if tt.wantErrorContains != "" {
				assert.ErrorContains(t, err, tt.wantErrorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantIdentity, resources)
			}
		})
	}
}

func TestApiClient_GetResources(t *testing.T) {
	token := "some-token"
	timeout := 5

	tests := []struct {
		name              string
		givenStatusCode   int
		givenData         []byte
		wantResources     ResourcesDTO
		wantErrorContains string
	}{
		{
			name:            "successfully fetches resources",
			givenStatusCode: 200,
			givenData: []byte(`
				<?xml version="1.0" encoding="UTF-8"?>\n
				<MediaContainer size="2">\n
					<Device name="MyPlexServer" product="Plex Media Server" productVersion="some-version" platform="some-platform" platformVersion="some-platform-version" device="some-device" clientIdentifier="1142ed040a27acc36ea876e8362b28464c3d240d" createdAt="1540597578" lastSeenAt="1681654399" provides="server" owned="1" accessToken="some-token" publicAddress="some-public-ip" httpsRequired="0" synced="0" relay="1" dnsRebindingProtection="0" natLoopbackSupported="0" publicAddressMatches="1" presence="1">\n
						<Connection protocol="https" address="some.private.i" port="32400" uri="https://some-private-ip.some-server-id.plex.direct:32400" local="1"/>\n
						<Connection protocol="https" address="some.public.ip" port="32400" uri="https://some-public-ip.some-server-id.plex.direct:32400" local="0"/>\n
					</Device>\n
					<Device name="OtherPlexServer" product="Plex Media Server" productVersion="some-version" platform="some-platform" platformVersion="some-platform-version" device="some-device" clientIdentifier="8cc5e1ff10756c3c4c7d0ada6189eabd06302cff" createdAt="1552689146" lastSeenAt="1681524851" provides="server" owned="0" publicAddress="some.public.ip" httpsRequired="0" ownerId="some-owner-id" home="0" accessToken="some-token" sourceTitle="some-source-title" synced="0" relay="1" dnsRebindingProtection="0" natLoopbackSupported="0" publicAddressMatches="0" presence="0">\n
						<Connection protocol="https" address="some.private.ipv4" port="32400" uri="https://some-private-ipv4.other-server-id.plex.direct:32400" local="1"/>\n
						<Connection protocol="https" address="some.public.ipv6" port="32400" uri="https://some-public-ipv6.other-server-id.plex.direct:32400" local="0"/>\n
					</Device>\n
				</MediaContainer>
			`),
			wantResources: ResourcesDTO{
				Devices: []DeviceDTO{
					{
						Name:             "MyPlexServer",
						Product:          "Plex Media Server",
						ClientIdentifier: "1142ed040a27acc36ea876e8362b28464c3d240d",
						Connections: []ConnectionDTO{
							{
								Protocol: "https",
								Address:  "some.private.i",
								URI:      "https://some-private-ip.some-server-id.plex.direct:32400",
								Local:    "1",
							},
							{
								Protocol: "https",
								Address:  "some.public.ip",
								URI:      "https://some-public-ip.some-server-id.plex.direct:32400",
								Local:    "0",
							},
						},
					},
					{
						Name:             "OtherPlexServer",
						Product:          "Plex Media Server",
						ClientIdentifier: "8cc5e1ff10756c3c4c7d0ada6189eabd06302cff",
						Connections: []ConnectionDTO{
							{
								Protocol: "https",
								Address:  "some.private.ipv4",
								URI:      "https://some-private-ipv4.other-server-id.plex.direct:32400",
								Local:    "1",
							},
							{
								Protocol: "https",
								Address:  "some.public.ipv6",
								URI:      "https://some-public-ipv6.other-server-id.plex.direct:32400",
								Local:    "0",
							},
						},
					},
				},
			},
		},
		{
			name:            "handles MediaContainer without devices",
			givenStatusCode: 200,
			givenData: []byte(`
				<?xml version="1.0" encoding="UTF-8"?>\n
				<MediaContainer size="0">\n
				</MediaContainer>
			`),
			wantResources: ResourcesDTO{},
		},
		{
			name:            "handles Device without connections",
			givenStatusCode: 200,
			givenData: []byte(`
				<?xml version="1.0" encoding="UTF-8"?>\n
				<MediaContainer size="1">\n
					<Device name="MyPlexServer" product="Plex Media Server" productVersion="some-version" platform="some-platform" platformVersion="some-platform-version" device="some-device" clientIdentifier="1142ed040a27acc36ea876e8362b28464c3d240d" createdAt="1540597578" lastSeenAt="1681654399" provides="server" owned="1" accessToken="some-token" publicAddress="some-public-ip" httpsRequired="0" synced="0" relay="1" dnsRebindingProtection="0" natLoopbackSupported="0" publicAddressMatches="1" presence="1">\n
					</Device>\n
				</MediaContainer>
			`),
			wantResources: ResourcesDTO{
				Devices: []DeviceDTO{
					{
						Name:             "MyPlexServer",
						Product:          "Plex Media Server",
						ClientIdentifier: "1142ed040a27acc36ea876e8362b28464c3d240d",
					},
				},
			},
		},
		{
			name:              "returns error for non-200 response code",
			givenStatusCode:   401,
			wantErrorContains: "failed with status code 401 (401 Unauthorized)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, token, r.Header.Get(headerKeyToken))
				assert.Equal(t, resourcesEndpoint, r.URL.Path)
				assert.Equal(t, plexTrue, r.URL.Query().Get(queryKeyIncludeHttps))
				assert.Equal(t, plexTrue, r.URL.Query().Get(queryKeyIncludeIPv6))

				w.WriteHeader(tt.givenStatusCode)
				_, err := w.Write(tt.givenData)
				require.NoError(t, err)
			}))

			client := NewApiClient(server.URL, token, timeout)

			// WHEN
			resources, err := client.GetResources()

			// THEN
			if tt.wantErrorContains != "" {
				assert.ErrorContains(t, err, tt.wantErrorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResources, resources)
			}
		})
	}
}

func TestApiClient_GetPreferences(t *testing.T) {
	token := "some-token"
	timeout := 5

	tests := []struct {
		name              string
		givenStatusCode   int
		givenData         []byte
		wantPreferences   PreferencesDTO
		wantErrorContains string
	}{
		{
			name:            "successfully fetches preferences",
			givenStatusCode: 200,
			givenData: []byte(`
				<?xml version="1.0" encoding="UTF-8"?>\n
				<MediaContainer size="3">\n
					<Setting id="FriendlyName" label="Friendly name" summary="This name will be used to identify this media server to other computers on your network. If you leave it blank, your computer&#39;s name will be used instead." type="text" default="" value="MyPlexServer" hidden="0" advanced="0" group="general" />\n
					<Setting id="sendCrashReports" label="Send crash reports to Plex" summary="This helps us improve your experience." type="bool" default="1" value="0" hidden="0" advanced="0" group="general" />\n
					<Setting id="customConnections" label="Custom server access URLs" summary="A comma-separated list of URLs (http or https) which are published up to plex.tv for server discovery." type="text" default="" value="" hidden="0" advanced="1" group="network" />\n
				</MediaContainer>
			`),
			wantPreferences: PreferencesDTO{
				Settings: []SettingDTO{
					{
						ID:    "FriendlyName",
						Type:  "text",
						Value: "MyPlexServer",
					},
					{
						ID:    "sendCrashReports",
						Type:  "bool",
						Value: "0",
					},
					{
						ID:    "customConnections",
						Type:  "text",
						Value: "",
					},
				},
			},
		},
		{
			name:            "handles MediaContainer without settings",
			givenStatusCode: 200,
			givenData: []byte(`
				<?xml version="1.0" encoding="UTF-8"?>\n
				<MediaContainer size="0">\n
				</MediaContainer>
			`),
			wantPreferences: PreferencesDTO{},
		},
		{
			name:              "returns error for non-200 response code",
			givenStatusCode:   401,
			wantErrorContains: "failed with status code 401 (401 Unauthorized)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, token, r.Header.Get(headerKeyToken))
				assert.Equal(t, preferencesEndpoint, r.URL.Path)

				w.WriteHeader(tt.givenStatusCode)
				_, err := w.Write(tt.givenData)
				require.NoError(t, err)
			}))

			client := NewApiClient(server.URL, token, timeout)

			// WHEN
			resources, err := client.GetPreferences()

			// THEN
			if tt.wantErrorContains != "" {
				assert.ErrorContains(t, err, tt.wantErrorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPreferences, resources)
			}
		})
	}
}

func TestApiClient_UpdateCustomConnections(t *testing.T) {
	token := "some-token"
	timeout := 5
	customConnections := "http://unused:32400/"

	tests := []struct {
		name              string
		givenStatusCode   int
		wantErrorContains string
	}{
		{
			name:            "successfully updates custom connections",
			givenStatusCode: 200,
		},
		{
			name:              "returns error for non-200 response code",
			givenStatusCode:   401,
			wantErrorContains: "failed with status code 401 (401 Unauthorized)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, token, r.Header.Get(headerKeyToken))
				assert.Equal(t, preferencesEndpoint, r.URL.Path)
				assert.Equal(t, customConnections, r.URL.Query().Get(queryKeyCustomConnections))

				w.WriteHeader(tt.givenStatusCode)
			}))

			client := NewApiClient(server.URL, token, timeout)

			// WHEN
			err := client.UpdateCustomConnections(customConnections)

			// THEN
			if tt.wantErrorContains != "" {
				assert.ErrorContains(t, err, tt.wantErrorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
