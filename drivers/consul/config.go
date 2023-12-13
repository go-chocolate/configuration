package consul

import (
	"net/http"
	"os"
	"time"

	capi "github.com/hashicorp/consul/api"
)

// Config
// consul 本身的连接配置
type Config struct {
	// Address is the address of the Consul server
	Address string

	// Scheme is the URI scheme for the Consul server
	Scheme string

	// Prefix for URIs for when consul is behind an API gateway (reverse
	// proxy).  The API gateway must strip off the PathPrefix before
	// passing the request onto consul.
	PathPrefix string

	// Datacenter to use. If not provided, the default agent datacenter is used.
	Datacenter string

	// Transport is the Transport to use for the http client.
	//Transport *http.Transport

	// HttpClient is the client to use. Default will be
	// used if not provided.
	//HttpClient *http.Client

	// HttpAuth is the auth info to use for http access.
	HttpAuth *capi.HttpBasicAuth

	// WaitTime limits how long a Watch will block. If not provided,
	// the agent default values will be used.
	WaitTime string

	// Token is used to provide a per-request ACL token
	// which overrides the agent's default token.
	Token string

	// TokenFile is a file containing the current token to use for this client.
	// If provided it is read once at startup and never again.
	TokenFile string

	// Namespace is the name of the namespace to send along for the request
	// when no other Namespace is present in the QueryOptions
	Namespace string

	// Partition is the name of the partition to send along for the request
	// when no other Partition is present in the QueryOptions
	Partition string

	transport  *http.Transport
	httpClient *http.Client
}

func (c *Config) consul() *capi.Config {
	return &capi.Config{
		Address:    c.Address,
		Scheme:     c.Scheme,
		PathPrefix: c.PathPrefix,
		Datacenter: c.Datacenter,
		HttpAuth:   c.HttpAuth,
		WaitTime:   duration(c.WaitTime),
		Token:      c.Token,
		TokenFile:  c.TokenFile,
		Namespace:  c.Namespace,
		Partition:  c.Partition,

		Transport:  c.transport,
		HttpClient: c.httpClient,
	}
}

func (c *Config) WithTransport(transport *http.Transport) {
	c.transport = transport
}

func (c *Config) WithHttpClient(client *http.Client) {
	c.httpClient = client
}

func duration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

func DefaultConfig() Config {
	address := os.Getenv("CONSUL_ADDRESS")
	if address == "" {
		address = "consul:8500"
	}
	config := Config{
		Address:    address,
		Scheme:     os.Getenv("CONSUL_SCHEME"),
		PathPrefix: os.Getenv("CONSUL_PATH_PREFIX"),
		Datacenter: os.Getenv("CONSUL_DATACENTER"),
		Token:      os.Getenv("CONSUL_TOKEN"),
		TokenFile:  os.Getenv("CONSUL_TOKEN_FILE"),
		Namespace:  os.Getenv("CONSUL_NAMESPACE"),
		Partition:  os.Getenv("CONSUL_PARTITION"),
	}

	username := os.Getenv("CONSUL_USERNAME")
	password := os.Getenv("CONSUL_PASSWORD")
	if username != "" {
		config.HttpAuth = &capi.HttpBasicAuth{
			Username: username,
			Password: password,
		}
	}
	return config
}
