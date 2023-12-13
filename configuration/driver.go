package configuration

import (
	"context"
	"os"
	"strings"
)

type Provider interface {
	GetContent(ctx context.Context, key string) ([]byte, string, error)
}

type ProviderFunc func(ctx context.Context, key string) ([]byte, string, error)

func (f ProviderFunc) GetContent(ctx context.Context, key string) ([]byte, string, error) {
	return f(ctx, key)
}

var LocalFileProvider ProviderFunc = func(ctx context.Context, key string) ([]byte, string, error) {
	var ext string
	if n := strings.LastIndex(key, "."); n > 0 {
		ext = key[n+1:]
	}
	content, err := os.ReadFile(key)
	if err != nil {
		return nil, ext, err
	}
	return content, ext, nil
}

func BytesProvider(data []byte, format string) ProviderFunc {
	return func(ctx context.Context, key string) ([]byte, string, error) {
		return data, format, nil
	}
}

func LocalFileDriver(script string) (Provider, error) {
	return LocalFileProvider, nil
}

type Driver func(script string) (Provider, error)

var drivers = map[string]Driver{
	"":      LocalFileDriver,
	"local": LocalFileDriver,
}

func RegisterConfigDriver(name string, driver Driver) {
	drivers[name] = driver
}
