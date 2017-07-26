package qtypes

import (
	"fmt"
	"strings"
	"sort"
	"bytes"
)

type Dimensions struct {
	Map map[string]string
}

func NewDimensions() Dimensions {
	return NewDimensionsPre(map[string]string{})
}

func NewDimensionsPre(dims map[string]string) Dimensions {
	return Dimensions{
		Map: dims,
	}
}

func NewDimensionsFromString(str string) Dimensions {
	dims := NewDimensions()
	for _, tupel := range strings.Split(str, ",") {
		kv := strings.Split(tupel,"=")
		if len(kv) != 2 {
			return dims
		} else {
			dims.Add(kv[0], kv[1])
		}
	}
	return dims
}

func NewDimensionsFromBytes(inp []byte) Dimensions {
	dims := NewDimensions()
	for _, tupel := range  bytes.SplitN(inp, []byte{','}, -1) {
		kv := bytes.SplitN(tupel,[]byte{'='},2)
		if len(kv) != 2 {
			return dims
		} else {
			dims.Add(string(kv[0]), string(kv[1]))
		}
	}
	return dims
}

func (dim *Dimensions) GetDims() map[string]string {
	return dim.Map
}


func (dim *Dimensions) Add(key,val string) {
	dim.Map[key] = val
}

func (dim *Dimensions) String() string {
	res := []string{}
		for k,v := range dim.Map {
		res = append(res, fmt.Sprintf("%s=%s", k,v))
	}
	sort.Strings(res)
	return strings.Join(res, ",")
}