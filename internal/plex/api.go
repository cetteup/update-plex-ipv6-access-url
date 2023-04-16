package plex

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	BaseURL                   = "https://plex.tv/api"
	resourcesEndpoint         = "/resources"
	preferencesEndpoint       = "/:/prefs"
	headerKeyToken            = "X-Plex-Token"
	queryKeyIncludeHttps      = "includeHttps"
	queryKeyIncludeIPv6       = "includeIPv6"
	queryKeyCustomConnections = "customConnections"
)

type GetResourcesDTO struct {
	Devices []DeviceDTO `xml:"Device"`
}

func (r GetResourcesDTO) GetDeviceByIdentifier(clientIdentifier string) (DeviceDTO, error) {
	for _, d := range r.Devices {
		if d.ClientIdentifier == clientIdentifier {
			return d, nil
		}
	}

	return DeviceDTO{}, fmt.Errorf("no such device: %s", clientIdentifier)
}

type DeviceDTO struct {
	Name             string          `xml:"name,attr"`
	Product          string          `xml:"product,attr"`
	ClientIdentifier string          `xml:"clientIdentifier,attr"`
	Connections      []ConnectionDTO `xml:"Connection"`
}

func (d DeviceDTO) GetLocalConnection() (ConnectionDTO, error) {
	return d.getConnectionByLocation(plexTrue)
}

func (d DeviceDTO) GetRemoteConnection() (ConnectionDTO, error) {
	return d.getConnectionByLocation(plexFalse)
}

func (d DeviceDTO) getConnectionByLocation(local string) (ConnectionDTO, error) {
	for _, c := range d.Connections {
		if c.Local == local {
			return c, nil
		}
	}

	return ConnectionDTO{}, fmt.Errorf("no location connection found for device: %s", d.Name)
}

func (d DeviceDTO) GetPlexDirectHostname() (string, error) {
	for _, c := range d.Connections {
		if strings.Contains(c.URI, ".plex.direct") {
			u, err := url.Parse(c.URI)
			if err != nil {
				return "", err
			}

			// Remove `dashed-ipv6-address.` prefix to only return `[server-id].plex.direct`
			hostname := strings.TrimPrefix(u.Hostname(), fmt.Sprintf("%s.", strings.ReplaceAll(c.Address, ".", "-")))
			return hostname, nil
		}
	}

	return "", fmt.Errorf("no .plex.direct hostname found for device: %s", d.Name)
}

type ConnectionDTO struct {
	Protocol string `xml:"protocol,attr"`
	Address  string `xml:"address,attr"`
	URI      string `xml:"uri,attr"`
	Local    string `xml:"local,attr"`
}

type ApiClient struct {
	client  http.Client
	baseURL string
	token   string
}

func NewApiClient(baseURL string, token string) *ApiClient {
	return &ApiClient{
		client: http.Client{
			Timeout: time.Second * 5,
		},
		baseURL: baseURL,
		token:   token,
	}
}

func (c *ApiClient) GetResources() (GetResourcesDTO, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return GetResourcesDTO{}, err
	}

	u = u.JoinPath(resourcesEndpoint)

	q := u.Query()
	q.Set(queryKeyIncludeHttps, plexTrue)
	q.Set(queryKeyIncludeIPv6, plexTrue)
	u.RawQuery = q.Encode()

	req, err := c.createRequest(http.MethodGet, u.String())
	if err != nil {
		return GetResourcesDTO{}, err
	}

	bytes, err := c.do(req)
	if err != nil {
		return GetResourcesDTO{}, err
	}

	var resources GetResourcesDTO
	if err := xml.Unmarshal(bytes, &resources); err != nil {
		return GetResourcesDTO{}, err
	}

	return resources, nil
}

func (c *ApiClient) UpdateCustomConnections(customConnections string) error {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}

	u = u.JoinPath(preferencesEndpoint)

	q := u.Query()
	q.Set(queryKeyCustomConnections, customConnections)
	u.RawQuery = q.Encode()

	req, err := c.createRequest(http.MethodPut, u.String())
	if err != nil {
		return err
	}

	_, err = c.do(req)
	return err
}

func (c *ApiClient) createRequest(method string, u string) (*http.Request, error) {
	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(headerKeyToken, c.token)

	return req, nil
}

func (c *ApiClient) do(req *http.Request) ([]byte, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err2 := res.Body.Close()
		if err2 != nil {
			log.Error().Err(err2).Msg("Failed to close Plex API request body")
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to %s failed with status code %d (%s)", res.Request.URL.String(), res.StatusCode, res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
