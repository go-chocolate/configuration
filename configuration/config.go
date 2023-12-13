package configuration

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/creasty/defaults"
	"github.com/ghodss/yaml"
)

var (
	configDriver = flag.String("driver", "local", "配置驱动，默认加载本地文件，可选值：local(本地文件)，consul，nacos")
	configScript = flag.String("script", "", "配置驱动加载脚本（consul 或 nacos 的连接配置）")
	configKey    = flag.String("config", "etc/conf.yaml", "配置路径")
)

func Load(v any) error {
	if !flag.Parsed() {
		flag.Parse()
	}

	_ = defaults.Set(v)

	driver, ok := drivers[*configDriver]
	if !ok {
		driver = LocalFileDriver
	}
	loader, err := driver(*configScript)
	if err != nil {
		return fmt.Errorf("load config driver failed: %v", err)
	}
	content, contentType, err := loader.GetContent(context.Background(), *configKey)
	if err != nil {
		return fmt.Errorf("load config content failed: %v", err)
	}
	content = []byte(ExpandTemplate(string(content), func(key string) string {
		var indent = 0
		if n := strings.LastIndex(key, ":"); n > 0 {
			indent, _ = strconv.Atoi(key[n+1:])
			key = key[:n]
		}
		val, _, _ := loader.GetContent(context.Background(), key)
		if len(val) == 0 || indent == 0 {
			return string(val)
		}
		var buf = bytes.NewBuffer([]byte{})
		for _, line := range bytes.Split(val, []byte{'\n'}) {
			buf.Write(bytes.Repeat([]byte{' '}, indent))
			buf.Write(line)
			buf.Write([]byte{'\n'})
		}
		return buf.String()
	}))
	content = []byte(ExpandEnv(string(content)))

	unmarshaler, ok := unmarshalers[strings.ToLower(contentType)]
	if !ok {
		unmarshaler = yaml.Unmarshal
	}
	if err := unmarshaler(content, v); err != nil {
		return fmt.Errorf("unmarshal config content failed: %v", err)
	}
	if validator, ok := v.(Validator); ok {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("validate config failed: %v", err)
		}
	}
	return nil
}

func MustLoad(v any) {
	if err := Load(v); err != nil {
		panic(err)
	}
}
