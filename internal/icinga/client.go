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
	defaultDialTimeout         = 10 * time.Second
	defaultKeepAlive           = 10 * time.Second
	defaultTLSHandshake        = 10 * time.Second
	defaultIdleConnTimeout     = 90 * time.Second
	defaultMaxIdleConns        = 100
	defaultMaxIdleConnsPerHost = 10
)

const (
	EndpointApiListener             = "/status/ApiListener"
	EndpointApplication             = "/status/IcingaApplication"
	EndpointCIB                     = "/status/CIB"
	EndpointCheckerComponent        = "/status/CheckerComponent"
	EndpointCompatLogger            = "/status/CompatLogger"
	EndpointElasticsearchWriter     = "/status/ElasticsearchWriter"
	EndpointExternalCommandListener = "/status/ExternalCommandListener"
	EndpointFileLogger              = "/status/FileLogger"
	EndpointGELFWriter              = "/status/GelfWriter"
	EndpointGraphiteWriter          = "/status/GraphiteWriter"
	EndpointIdoMysqlConnection      = "/status/IdoMysqlConnection"
	EndpointIdoPgsqlConnection      = "/status/IdoPgsqlConnection"
	EndpointInfluxdb2Writer         = "/status/Influxdb2Writer"
	EndpointInfluxdbWriter          = "/status/InfluxdbWriter"
	EndpointJournaldLogger          = "/status/JournaldLogger"
	EndpointLivestatusListener      = "/status/LivestatusListener"
	EndpointNotificationComponent   = "/status/NotificationComponent"
	EndpointOpenTsdbWriter          = "/status/OpenTsdbWriter"
	EndpointOTLPMetricsWriter       = "/status/OTLPMetricsWriter"
	EndpointPerfdataWriter          = "/status/PerfdataWriter"
	EndpointSyslogLogger            = "/status/SyslogLogger"
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

// IcingaClient is an interface that we use to simplify testing
// Note, the methods use a context internally.
type IcingaClient interface {
	GetPerfdataMetrics(endpoint string) (PerfdataResult, error)
	GetCIBMetrics() (CIBResult, error)
	GetApplicationMetrics() (ApplicationResult, error)
	GetAPIListenerMetrics() (APIResult, error)
	GetCheckerComponentMetrics() (CheckerComponentResult, error)
}

// NewClient returns a configured client
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
			Timeout:   defaultDialTimeout,
			KeepAlive: defaultKeepAlive,
		}).DialContext,
		TLSHandshakeTimeout: defaultTLSHandshake,
		TLSClientConfig:     tlsConfig,
		IdleConnTimeout:     defaultIdleConnTimeout,
		MaxIdleConns:        defaultMaxIdleConns,
		MaxIdleConnsPerHost: defaultMaxIdleConnsPerHost,
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

// GetPerfdataMetrics returns the perfdata from a given status API endpoint
// There is some duplication here, but that is fine for now
func (c *Client) GetPerfdataMetrics(endpoint string) (PerfdataResult, error) {
	var result PerfdataResult

	body, errBody := c.fetchJSON(endpoint)

	if errBody != nil {
		return result, fmt.Errorf("error fetching response: %w", errBody)
	}

	errDecode := json.Unmarshal(body, &result)

	if errDecode != nil {
		return result, fmt.Errorf("error parsing response: %w", errDecode)
	}

	if len(result.Results) < 1 {
		return result, fmt.Errorf("no results for '%s' endpoint", endpoint)
	}

	return result, nil
}

// GetCIBMetrics returns the Common Information Base metrics
func (c *Client) GetCIBMetrics() (CIBResult, error) {
	var result CIBResult

	body, errBody := c.fetchJSON(EndpointCIB)

	if errBody != nil {
		return result, fmt.Errorf("error fetching response: %w", errBody)
	}

	errDecode := json.Unmarshal(body, &result)

	if errDecode != nil {
		return result, fmt.Errorf("error parsing response: %w", errDecode)
	}

	return result, nil
}

// GetApplicationMetrics returns the base application metrics
func (c *Client) GetApplicationMetrics() (ApplicationResult, error) {
	var result ApplicationResult

	body, errBody := c.fetchJSON(EndpointApplication)

	if errBody != nil {
		return result, fmt.Errorf("error fetching response: %w", errBody)
	}

	errDecode := json.Unmarshal(body, &result)

	if errDecode != nil {
		return result, fmt.Errorf("error parsing response: %w", errDecode)
	}

	return result, nil
}

// GetAPIListenerMetrics returns the APIListener metrics
func (c *Client) GetAPIListenerMetrics() (APIResult, error) {
	var result APIResult

	body, errBody := c.fetchJSON(EndpointApiListener)

	if errBody != nil {
		return result, fmt.Errorf("error fetching response: %w", errBody)
	}

	errDecode := json.Unmarshal(body, &result)

	if errDecode != nil {
		return result, fmt.Errorf("error parsing response: %w", errDecode)
	}

	return result, nil
}

// GetCheckerComponentMetrics returns the CheckerComponent metrics
func (c *Client) GetCheckerComponentMetrics() (CheckerComponentResult, error) {
	var result CheckerComponentResult

	body, errBody := c.fetchJSON(EndpointCheckerComponent)

	if errBody != nil {
		return result, fmt.Errorf("error fetching response: %w", errBody)
	}

	errDecode := json.Unmarshal(body, &result)

	if errDecode != nil {
		return result, fmt.Errorf("error parsing response: %w", errDecode)
	}

	return result, nil
}

// fetchJSON calls the given endpoint and returns the JSON result
func (c *Client) fetchJSON(endpoint string) ([]byte, error) {
	// Lookup data in the cache we go out and bother the Icinga API
	if elem, ok := c.cache.Get(endpoint); ok {
		return elem, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u := c.URL.JoinPath(endpoint)

	req, errReq := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)

	if errReq != nil {
		return []byte{}, fmt.Errorf("error creating request: %w", errReq)
	}

	resp, errDo := c.Client.Do(req)

	if errDo != nil {
		return []byte{}, fmt.Errorf("error performing request: %w", errDo)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		//nolint: errcheck
		io.Copy(io.Discard, resp.Body)
		return []byte{}, fmt.Errorf("request failed: %s", resp.Status)
	}

	data, errRead := io.ReadAll(resp.Body)

	if errRead != nil {
		return []byte{}, fmt.Errorf("reading response failed: %w", errRead)
	}

	c.cache.Set(endpoint, data)

	return data, nil
}
