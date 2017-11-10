package qtypes_helper

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCompareMap(t *testing.T) {
	m := map[string]interface{}{
		"name": "name",
		"list": []string{"1","2"},
		"map": map[string]string{"key": "val"},
		"default": []float64{1.2,1.3},
	}
	assert.True(t, CompareMap(m,m))
	fail1 := map[string]interface{}{
		"list": []string{"1","2"},
		"map": map[string]string{"key": "val"},
	}
	assert.False(t, CompareMap(m,fail1))
	fail2 := map[string]interface{}{
		"name": "name1",
		"list": []string{"1","2"},
		"map": map[string]string{"key": "val"},
		"default": []float64{1.2,1.3},
	}
	assert.False(t, CompareMap(m,fail2))
	fail3 := map[string]interface{}{
		"name": "name",
		"list": []string{"1","2","3"},
		"map": map[string]string{"key": "val"},
		"default": []float64{1.2,1.3},
	}
	assert.False(t, CompareMap(m,fail3))
	fail4 := map[string]interface{}{
		"name": "name",
		"list": []string{"1","2"},
		"map": map[string]string{"key": "val1"},
		"default": []float64{1.2,1.3},
	}
	assert.False(t, CompareMap(m,fail4))
	fail5 := map[string]interface{}{
		"name": "name",
		"list": []string{"1","2"},
		"map": map[string]string{"key": "val"},
		"default": []float64{1.2,2.4},
	}
	assert.False(t, CompareMap(m,fail5))
}

func TestPrefixFlatKV(t *testing.T) {
	kv := map[string]string{
		"key1": "val1",
		"key2": "val2",
	}
	out := map[string]interface{}{
		"org1": "val1",
		"org2": "val2",
	}
	exp := map[string]interface{}{
		"org1": "val1",
		"org2": "val2",
		"msg_key1": "val1",
		"msg_key2": "val2",
	}
	got, err := PrefixFlatKV(kv, out, "msg")
	assert.NoError(t, err)
	assert.Equal(t, exp, got)
	kv["key 3"] = "val3"
	_, err = PrefixFlatKV(kv, out, "msg")
	assert.Error(t, err)
}

func TestPrefixFlatLabels(t *testing.T) {
	lab := []string{
		"key1=val1",
		"key2=val2",
	}
	out := map[string]interface{}{
		"org1": "val1",
		"org2": "val2",
	}
	exp := map[string]interface{}{
		"org1": "val1",
		"org2": "val2",
		"msg_key1": "val1",
		"msg_key2": "val2",
	}
	got, err := PrefixFlatLabels(lab, out, "msg")
	assert.NoError(t, err)
	assert.Equal(t, exp, got)
	lab = append(lab, "key 3=val3")
	_, err = PrefixFlatLabels(lab, out, "msg")
	assert.Error(t, err)
}
func TestFlattenKV(t *testing.T) {
	kv := map[string]string{
		"key1": "val1",
		"key2": "val2",
	}
	got, err := FlattenKV(kv)
	assert.NoError(t, err)
	assert.Equal(t, "key1=val1,key2=val2", got)
	kv["key3"] = "hello,world"
	_, err = FlattenKV(kv)
	assert.Error(t, err)
	kv["key3"] = "HelloWorld"
	kv["key 4"] = "fail"
	_, err = FlattenKV(kv)
	assert.Error(t, err)
	delete(kv, "key 4")
	_, err = FlattenKV(kv)
	assert.NoError(t, err)
	kv["key,4"] = "fail"
	_, err = FlattenKV(kv)
	assert.Error(t, err)
}
