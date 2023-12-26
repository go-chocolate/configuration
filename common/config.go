package common

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Config map[string]any

func (c Config) Get(key string) any {
	return c[key]
}

func (c Config) String(key string) string {
	val, ok := c[key]
	if !ok {
		return ""
	}
	return any2string(val)
}

func (c Config) Int64(key string) int64 {
	val, ok := c[key]
	if !ok {
		return 0
	}
	return any2int64(val)
}

func (c Config) Bool(key string) bool {
	v, ok := c[key]
	if !ok {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val == "true"
	}
	return false
}

func (c Config) Strings(key string) []string {
	v := c[key]
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case string:
		return []string{val}
	case []string:
		return val
	case []any:
		var result = make([]string, 0, len(val))
		for _, item := range val {
			result = append(result, any2string(item))
		}
		return result
	}
	return nil
}

func (c Config) Int64s(key string) []int64 {
	v := c[key]
	if v == nil {
		return nil
	}
	if val, ok := v.([]int64); ok {
		return val
	}
	switch val := v.(type) {
	case int64:
		return []int64{val}
	case []int64:
		return val
	case []any:
		var result = make([]int64, 0, len(val))
		for _, item := range val {
			result = append(result, any2int64(item))
		}
		return result
	case []json.Number:
		var result = make([]int64, 0, len(val))
		for _, item := range val {
			i, _ := item.Int64()
			result = append(result, i)
		}
		return result
	}
	return nil
}

func (c Config) Duration(key string) time.Duration {
	val := c.String(key)
	if isNumber([]byte(val)) {
		d, _ := strconv.Atoi(val)
		return time.Duration(d) * time.Millisecond
	}
	dur, _ := time.ParseDuration(val)
	return dur
}

func any2string(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func any2int64(v any) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case json.Number:
		i, _ := val.Int64()
		return i
	case string:
		i, _ := strconv.ParseInt(val, 10, 64)
		return i
	}
	return 0
}

func isNumber(data []byte) bool {
	var dot bool
	for _, v := range data {
		if v == '.' {
			if !dot {
				dot = true
				continue
			} else {
				return false
			}
		}
		if v < '0' || v > '9' {
			return false
		}
	}
	return true
}
