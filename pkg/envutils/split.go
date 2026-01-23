package envutils

import "strings"

func SplitEnv(key, fallback string) []string {
	value := GetEnv(key, fallback)

	if value == "" {
		return []string{}
	}

	return strings.Split(value, ",")
}
