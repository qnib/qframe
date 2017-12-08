package qtypes_helper

import (
	"fmt"
	"reflect"
	"strings"
	"sort"
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

// PrefixFlatKV takes a key/value map and merges each key into an existing map using a prefix.
func PrefixFlatKV(kv  map[string]string, out map[string]interface{}, prefix string) (res map[string]interface{},err error) {
	res = map[string]interface{}{}
	for k,v := range out {
		res[k] = v
	}
	forbiddenStr := []string{","," "}
	for k,v := range kv {
		if ContainersOneOf(k, forbiddenStr) {
			return res, fmt.Errorf("key containers one if '%s'", strings.Join(forbiddenStr, "','"))
		}
		res[fmt.Sprintf("%s_%s", prefix, k)] = v
	}
	return res,nil
}

// PrefixFlatLabels takes a key/value slice (separated by =) and merges each key into an existing map using a prefix.
func PrefixFlatLabels(s []string, out map[string]interface{}, prefix string) (res map[string]interface{},err error) {
	res = map[string]interface{}{}
	for k,v := range out {
		res[k] = v
	}
	forbiddenStr := []string{","," "}
	for _, ele := range s {
		kv := strings.Split(ele, "=")
		if len(kv) != 2 {
			continue
		}
		if ContainersOneOf(kv[0], forbiddenStr) {
			return res, fmt.Errorf("key containers one if '%s'", strings.Join(forbiddenStr, "','"))
		}
		res[fmt.Sprintf("%s_%s", prefix, kv[0])] = kv[1]
	}
	return res, nil
}




func FlattenKV(kv map[string]string) (res string, err error) {
	r := []string{}
	forbiddenStr := []string{","," "}
	for k,v := range kv {
		if ContainersOneOf(k, forbiddenStr) || ContainersOneOf(v, forbiddenStr) {
			return "", fmt.Errorf("key or value containers one if '%s'", strings.Join(forbiddenStr, "','"))
		}
		r = append(r, fmt.Sprintf("%s=%s", k,v))
	}
	sort.Strings(r)
	res = strings.Join(r, ",")
	return
}

func ContainersOneOf(s string, m []string) bool {
	for _, i := range m {
		if strings.Contains(s, i) {
			return true
		}
	}
	return false
}
