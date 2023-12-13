package configuration

import (
	"bytes"
	"os"
	"strings"
)

// ExpandEnv 环境变量占位符
// 格式： {NAME} (不带默认值)，{NAME:hello} (带默认值)
func ExpandEnv(s string) string {
	return expand(s, []byte{'$', '{'}, []byte{'}'}, GetEnv)
}

func ExpandTemplate(s string, get func(string) string) string {
	return expand(s, []byte("<<"), []byte(">>"), get)
}

func Expand(s string, left, right []byte, get func(string) string) string {
	return expand(s, left, right, get)
}

func expand(text string, left, right []byte, f func(string) string) string {

	var buf []byte
	var key []byte
	var expanding bool

	var i = 0
	for i < len(text) {
		if len(text) >= i+len(left) && bytes.Equal([]byte(text[i:i+len(left)]), left) {
			if expanding {
				buf = append(buf, left...)
			}
			buf = append(buf, key...)
			key = []byte{}
			expanding = true
			i += len(left)
			continue
		}
		if len(text) >= i+len(right) && expanding && bytes.Equal([]byte(text[i:i+len(right)]), right) {
			buf = append(buf, f(string(key))...)
			key = []byte{}
			expanding = false
			i += len(right)
			continue
		}
		if expanding {
			key = append(key, text[i])
		} else {
			buf = append(buf, text[i])
		}
		i++
	}
	if expanding {
		buf = append(buf, left...)
		buf = append(buf, key...)
	}
	return string(buf)
}

func GetEnv(key string) string {
	var k, def = key, ""
	if n := strings.Index(key, ":"); n > 0 {
		k = key[:n]
		def = key[n+1:]
	}
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
