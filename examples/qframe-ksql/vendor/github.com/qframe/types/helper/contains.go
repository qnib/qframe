package qtypes_helper

import (
	"strings"
	"fmt"
)


func IsInput(list []string, item string) bool {
	return IsItem(list, item)
}

// IsInput checks if a list of inputs (e.g. [one, two]) matches the suffix of item (two-one)
func IsItem(list []string, item string) bool {
	if item == "*" {
		return true
	}
	prefixItem := fmt.Sprintf("->%s", item)
	for _, i := range list {
		if item == i {
			return true
		}
		if strings.HasSuffix(prefixItem, i) {
			return true
		}
	}
	return false
}

// IsLastSource checks if a list of inputs (e.g. [one, two]) matches the last item in SourcesPath of item [two,one]
func IsLastSource(list []string, item string) bool {
	return list[len(list)-1] == item
}