package configuration

import (
	"encoding/json"

	"github.com/ghodss/yaml"
)

type Unmarshaler func(b []byte, v any) error

var unmarshalers = map[string]Unmarshaler{
	"json": json.Unmarshal,
	"yaml": yaml.Unmarshal,
	"":     yaml.Unmarshal,
}

func RegisterUnmarshaler(name string, unmarshaler Unmarshaler) {
	unmarshalers[name] = unmarshaler
}
