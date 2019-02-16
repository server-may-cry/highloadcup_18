package handlers

import (
	"bytes"
	"strconv"
)

func stringWithComasToInts(query []byte) []int {
	strings := stringWithComasToStrings(query)
	ints := make([]int, 0, len(strings))
	for _, v := range strings {
		integer, _ := strconv.Atoi(v)
		ints = append(ints, integer)
	}
	return ints
}

func stringWithComasToStrings(query []byte) []string {
	parts := bytes.Split(query, []byte(","))
	result := make([]string, 0, len(parts))
	for _, v := range parts {
		s := string(v)
		if s == "" {
			return []string{}
		}
		result = append(result, s)
	}
	return result
}
