package array

import "strings"

func ToString(x [65]int8) string {
	var buf [65]byte
	for i, b := range x {
		buf[i] = byte(b)
	}
	str := string(buf[:])
	if i := strings.Index(str, "\x00"); i != -1 {
		str = str[:i]
	}
	return str
}

func InArray(s interface{}, d []interface{}) bool {
	for _, v := range d {
		if s == v {
			return true
		}
	}
	return false
}

func ArrayColumn(d map[int]map[string]string, column_key string, index_key string) map[int]map[string]string {
	nd := make(map[int]map[string]string)
	for k, v := range d {
		for e, q := range v {
			if e == index_key {
				nd[k][index_key] = q
			}
		}
	}
	return nd
}
