package icinga

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	endpointApiListener             = "/status/ApiListener"
	endpointApplication             = "/status/IcingaApplication"
	endpointCIB                     = "/status/CIB"
	endpointCheckerComponent        = "/status/CheckerComponent"
	endpointCompatLogger            = "/status/CompatLogger"
	endpointElasticsearchWriter     = "/status/ElasticsearchWriter"
	endpointExternalCommandListener = "/status/ExternalCommandListener"
	endpointFileLogger              = "/status/FileLogger"
	endpointGelfWriter              = "/status/GelfWriter"
	endpointGraphiteWriter          = "/status/GraphiteWriter"
	endpointIcingaApplication       = "/status/IcingaApplication"
	endpointIdoMysqlConnection      = "/status/IdoMysqlConnection"
	endpointIdoPgsqlConnection      = "/status/IdoPgsqlConnection"
	endpointInfluxdb2Writer         = "/status/Influxdb2Writer"
	endpointInfluxdbWriter          = "/status/InfluxdbWriter"
	endpointJournaldLogger          = "/status/JournaldLogger"
	endpointLivestatusListener      = "/status/LivestatusListener"
	endpointNotificationComponent   = "/status/NotificationComponent"
	endpointOpenTsdbWriter          = "/status/OpenTsdbWriter"
	endpointPerfdataWriter          = "/status/PerfdataWriter"
	endpointSyslogLogger            = "/status/SyslogLogger"
)

type Config struct {
	BasicAuthUsername string
	BasicAuthPassword string
	CAFile            string
	CertFile          string
	KeyFile           string
	Insecure          bool
	CacheTTL          time.Duration
	IcingaAPIURI      url.URL
}

type Client struct {
	Client http.Client
	URL    url.URL
	cache  *Cache
	config Config
}

func NewClient(c Config) (*Client, error) {
	// Create TLS configuration for default RoundTripper
	tlsConfig, err := newTLSConfig(&TLSConfig{
		InsecureSkipVerify: c.Insecure,
		CAFile:             c.CAFile,
		KeyFile:            c.KeyFile,
		CertFile:           c.CertFile,
	})

	if err != nil {
		return nil, err
	}

	var rt http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     tlsConfig,
	}

	// Using a BasicAuth for authentication
	if c.BasicAuthUsername != "" {
		rt = newBasicAuthRoundTripper(c.BasicAuthUsername, c.BasicAuthPassword, rt)
	}

	cache := NewCache(c.CacheTTL)

	cli := &Client{
		URL: c.IcingaAPIURI,
		Client: http.Client{
			Transport: rt,
		},
		config: c,
		cache:  cache,
	}

	return cli, nil
}

func (icinga *Client) GetApiListenerMetrics() (APIResult, error) {
	var result APIResult

	body, errBody := icinga.fetchJSON(endpointApiListener)

	if errBody != nil {
		return result, fmt.Errorf("error fetching response: %w", errBody)
	}

	errDecode := json.Unmarshal(body, &result)

	if errDecode != nil {
		return result, fmt.Errorf("error parsing response: %w", errDecode)
	}

	return result, nil
}

func (icinga *Client) GetCIBMetrics() (CIBResult, error) {
	var result CIBResult

	body, errBody := icinga.fetchJSON(endpointCIB)

	if errBody != nil {
		return result, fmt.Errorf("error fetching response: %w", errBody)
	}

	errDecode := json.Unmarshal(body, &result)

	if errDecode != nil {
		return result, fmt.Errorf("error parsing response: %w", errDecode)
	}

	return result, nil
}

func (icinga *Client) GetApplicationMetrics() (ApplicationResult, error) {
	var result ApplicationResult

	body, errBody := icinga.fetchJSON(endpointApplication)

	if errBody != nil {
		return result, fmt.Errorf("error fetching response: %w", errBody)
	}

	errDecode := json.Unmarshal(body, &result)

	if errDecode != nil {
		return result, fmt.Errorf("error parsing response: %w", errDecode)
	}

	return result, nil
}

func (icinga *Client) GetCheckerComponentMetrics() (CheckerComponentResult, error) {
	var result CheckerComponentResult

	body, errBody := icinga.fetchJSON(endpointCheckerComponent)

	if errBody != nil {
		return result, fmt.Errorf("error fetching response: %w", errBody)
	}

	errDecode := json.Unmarshal(body, &result)

	if errDecode != nil {
		return result, fmt.Errorf("error parsing response: %w", errDecode)
	}

	return result, nil
}

func (icinga *Client) fetchJSON(endpoint string) ([]byte, error) {
	// Lookup data in the cache we go out and bother the Icinga API
	if elem, ok := icinga.cache.Get(endpoint); ok {
		return elem, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := icinga.URL.JoinPath(endpoint)

	req, errReq := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)

	if errReq != nil {
		return []byte{}, fmt.Errorf("error creating request: %w", errReq)
	}

	resp, errDo := icinga.Client.Do(req)

	if errDo != nil {
		return []byte{}, fmt.Errorf("error performing request: %w", errDo)
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("request failed: %s", resp.Status)
	}

	defer resp.Body.Close()

	data, errRead := io.ReadAll(resp.Body)

	if errRead != nil {
		return []byte{}, fmt.Errorf("reading response failed: %w", errRead)
	}

	icinga.cache.Set(endpoint, data)

	return data, nil
}
