package configuration

import (
	"flag"
	"os"
	"testing"
)

type Global struct {
	Name string ``
	Age  int    ``
}

type TestConfig struct {
	Name   string ``
	Score  int    `default:"100"`
	Host   string ``
	Global Global ``
}

func TestLoadConfig(t *testing.T) {
	os.Setenv("HOST", "127.0.0.1")
	flag.Set("config", "testdata/conf.yaml")
	var c = new(TestConfig)
	if err := Load(c); err != nil {
		t.Error(err)
	}
	t.Log(c)
	if c.Name != "Alex" || c.Score != 100 || c.Host != "127.0.0.1" || c.Global.Name != "Bob" || c.Global.Age != 32 {
		t.Fail()
	}
}

func TestLoadConfigWithCustomDriver(t *testing.T) {
	flag.Set("driver", "custom")

	RegisterConfigDriver("custom", func(script string) (Provider, error) {
		return BytesProvider([]byte("Name: Alen\nScore: 50"), "yaml"), nil
	})

	var c = new(TestConfig)
	if err := Load(c); err != nil {
		t.Error(err)
	}
	t.Log(c.Name, c.Score)
	if c.Name != "Alen" || c.Score != 50 {
		t.Fail()
	}
}
