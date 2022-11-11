package common

import "os"

func GetOrDefault(name string, defaultVal string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	return defaultVal
}
