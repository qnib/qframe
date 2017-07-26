package qutils

import "regexp"

func GetParams(regEx *regexp.Regexp, str string) (paramsMap map[string]string) {
	match := regEx.FindStringSubmatch(str)
	paramsMap = make(map[string]string)
	for i, name := range regEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}

