package env

import "os"

func unwrapString(s string) string {
	if len(s) < 2 {
		return s
	}
	if s[0] == '"' {
		s = s[1:]
	}
	if i := len(s) - 1; s[i] == '"' {
		s = s[:i]
	}
	return s
}

func Get(key string, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		val = fallback
	}
	return unwrapString(val)
}

func GetBool(key string, fallback bool) bool {
	val := fallback
	str, ok := os.LookupEnv(key)
	if ok {
		val = unwrapString(str) == "true"
	}
	return val
}
