package consul

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	capi "github.com/hashicorp/consul/api"
)

type ContentProvider struct {
	client *capi.Client
}

func NewContentProvider(config Config) (*ContentProvider, error) {
	client, err := capi.NewClient(config.consul())
	if err != nil {
		return nil, err
	}
	return &ContentProvider{client: client}, nil
}

func (p *ContentProvider) GetContent(ctx context.Context, key string) ([]byte, string, error) {
	kv, _, err := p.client.KV().Get(key, nil)
	if err != nil {
		return nil, "", err
	}
	//key not found
	if kv == nil {
		return []byte{}, "", nil
	}
	var ext string
	if n := strings.LastIndex(key, "."); n > 0 {
		ext = key[n+1:]
	}
	return kv.Value, ext, nil
}

func Driver(script string) (*ContentProvider, error) {
	content, err := os.ReadFile(script)
	if err != nil {
		return nil, err
	}

	var ext string
	if n := strings.LastIndex(script, "."); n > 0 {
		ext = strings.ToLower(script[n+1:])
	}
	var unmarshaler func([]byte, any) error
	switch ext {
	case "json":
		unmarshaler = json.Unmarshal
	default:
		unmarshaler = yaml.Unmarshal
	}
	var config Config
	if err := unmarshaler(content, &config); err != nil {
		return nil, err
	}

	return NewContentProvider(config)
}
