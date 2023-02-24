package utils

import "strconv"

func ParseIntDefault(s string, def int) int {
	return int(ParseInt64Default(s, int64(def)))
}

func ParseInt64Default(s string, def int64) int64 {
	l, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return def
	}

	return l
}
