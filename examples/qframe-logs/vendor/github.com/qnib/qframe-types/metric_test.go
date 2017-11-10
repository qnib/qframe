package qtypes


import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestNewMetric(t *testing.T) {
	m := New("testSrc", "testMetric")
	assert.Equal(t, "testMetric", m.Name)
	assert.Equal(t, "testSrc", m.SourcePath[0])
}

func TestNewExtMetric(t *testing.T) {
	now := time.Now()
	dims := map[string]string{
		"key1": "val1",
	}
	m := NewExt("testSrc", "testMetric", Gauge, 1.0, dims, now, false)
	assert.Equal(t, "gauge", m.MetricType)
	assert.Equal(t, 1.0, m.Value)
	assert.Equal(t, dims, m.Dimensions)
	assert.Equal(t, now, m.Time)
	assert.Equal(t, false, m.Buffered)
}

func TestMetric_GetDimensionList(t *testing.T) {
	now := time.Now()
	dims := map[string]string{
		"key1": "val1",
		"key2": "val2",
	}
	m := NewExt("testSrc", "testMetric", Gauge, 1.0, dims, now, false)
	got := m.GetDimensionList()
	assert.Equal(t, "key1=val1,key2=val2", got)
}