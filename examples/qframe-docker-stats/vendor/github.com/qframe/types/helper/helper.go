package qtypes_helper

import (
	"fmt"
	"reflect"
)

func CompareMap(exp, got map[string]interface{}) bool {
	for eK, eV := range exp {
		gV, ok := got[eK]
		if ! ok {
			fmt.Printf("Expected key '%s' not found\n", eK)
			return false
		}
		switch eV.(type) {
		case string,int,int64,float64,bool:
			if eV != gV {
				fmt.Printf("Key '%s' differs: expected:%v != %v\n", eK, eV, gV)
				return false
			}
		case []string:
			if ! reflect.DeepEqual(eV, gV) {
				fmt.Printf("Key '%s' differs: expected:%v != %v\n", eK, eV, gV)
				return false
			}
		case map[string]string:
			if ! reflect.DeepEqual(eV, gV) {
				fmt.Printf("Key '%s' differs: expected:%v != %v\n", eK, eV, gV)
				return false
			}
		default:
			if ! reflect.DeepEqual(eV, gV) {
				fmt.Printf("Key '%s' differs: expected:%v != %v\n", eK, eV, gV)
				return false
			}
		}

	}
	return true
}
