package qtypes_metrics

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
	"github.com/qframe/types/messages"
)


// The different types of metrics that are supported
const (
	Gauge             = "gauge"
	Counter           = "counter"
	CumulativeCounter = "cumcounter"
)

// Metric type holds all the information for a single metric data
// point. Metrics are generated in collectors and passed to handlers.
type Metric struct {
	qtypes_messages.Base
	Name       string            `json:"name"`
	MetricType string            `json:"type"`
	Value      float64           `json:"value"`
	Dimensions map[string]string `json:"dimensions"`
	Buffered   bool              `json:"buffered"`
}

// New returns a new metric with name. Default metric type is "gauge"
// and timestamp is set to now. Value is initialized to 0.0.
func New(source, name string) Metric {
	return Metric{
		Base: 		qtypes_messages.NewBase(source),
		Name:       sanitizeString(name),
		MetricType: Gauge,
		Value:      0.0,
		Dimensions: make(map[string]string),
		Buffered:   false,
	}
}

// NewExt provides a more controled creation
func NewExt(source, name string, metricTyp string, val float64, dimensions map[string]string, t time.Time, buffered bool) Metric {
	m := Metric{
		Base: 		qtypes_messages.NewTimedBase(source, t),
		Name:       sanitizeString(name),
		MetricType: metricTyp,
		Value:      val,
		Dimensions: dimensions,
		Buffered:   buffered,
	}
	return m
}

func (m *Metric) GetDimensionList() string {
	res := []string{}
	for k, v := range m.Dimensions {
		res = append(res, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(res, ",")
}

// Filter provides a struct that can filter a metric by Name (regex), type, dimension (subset of Dimensions)
type Filter struct {
	Name       string            `json:"name"`
	MetricType string            `json:"type"`
	Dimensions map[string]string `json:"dimensions"`
}

// ToJSON Transforms Filter to JSON
func (f *Filter) ToJSON() string {
	b, err := json.Marshal(f)
	if err != nil {
		fmt.Println(err)
		return "{}"
	}
	return string(b)
}

// NewFilter returns a Filter with compiled regex
func NewFilter(name string, t string, d map[string]string) Filter {
	return Filter{
		Name:       name,
		MetricType: t,
		Dimensions: d,
	}
}

// WithValue returns metric with value of type Gauge
func WithValue(source, name string, value float64) Metric {
	metric := New(source, name)
	metric.Value = value
	return metric
}

// EnableBuffering puts the metric into buffering handlers (e.g. ZmqBUF)
func (m *Metric) EnableBuffering() {
	m.Buffered = true
}

// DisableBuffering takes the metric out of buffering (e.g. ZmqBUF)
func (m *Metric) DisableBuffering() {
	m.Buffered = false
}

// SetTime to metric
func (m *Metric) SetTime(mtime time.Time) {
	m.Time = mtime
}

// AddDimension adds a new dimension to the Metric.
func (m *Metric) AddDimension(name, value string) {
	m.Dimensions[sanitizeString(name)] = sanitizeString(value)
}

// RemoveDimension removes a dimension from the Metric.
func (m *Metric) RemoveDimension(name string) {
	delete(m.Dimensions, name)
}

// AddDimensions adds multiple new dimensions to the Metric.
func (m *Metric) AddDimensions(dimensions map[string]string) {
	for k, v := range dimensions {
		m.AddDimension(k, v)
	}
}

// GetDimensions returns the dimensions of a metric merged with defaults. Defaults win.
func (m *Metric) GetDimensions(defaults map[string]string) (dimensions map[string]string) {
	dimensions = make(map[string]string)
	for name, value := range m.Dimensions {
		dimensions[name] = value
	}
	for name, value := range defaults {
		dimensions[name] = value
	}
	return dimensions
}

// GetDimensionValue returns the value of a dimension if it's set.
func (m *Metric) GetDimensionValue(dimension string) (value string, ok bool) {
	dimension = sanitizeString(dimension)
	value, ok = m.Dimensions[dimension]
	return
}

func (m *Metric) GetDimensionString() string {
	res := []string{}
	for k,v := range m.Dimensions {
		res = append(res, fmt.Sprintf("%s=%s", k,v))
	}
	return strings.Join(res, " ")
}
// AddToAll adds a map of dimensions to a list of metrics
func AddToAll(metrics *[]Metric, dims map[string]string) {
	for _, m := range *metrics {
		for key, value := range dims {
			m.AddDimension(key, value)
		}
	}
}

func sanitizeString(s string) string {
	s = strings.Replace(s, "=", "-", -1)
	s = strings.Replace(s, ":", "-", -1)
	return s
}

// ToJSON Transforms metric to JSON
func (m *Metric) ToJSON() string {
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return "{}"
	}
	return string(b)
}

// IsSubDim checks if the 2st map contains all items in the second
func (m *Metric) IsSubDim(other map[string]string) bool {
	for k, v := range other {
		val, ok := m.Dimensions[k]
		if !ok || v != val {
			return false
		}
	}
	return true
}

// IsFiltered checks if metrics is filtered with a given filter
func (m *Metric) IsFiltered(f Filter) bool {
	if !m.IsSubDim(f.Dimensions) {
		return false
	}
	if m.MetricType != f.MetricType {
		return false
	}
	// TODO: Precompile regex to speed up matching
	if !regexp.MustCompile(f.Name).MatchString(m.Name) {
		return false
	}

	return true
}

func (m *Metric) ToOpenTSDB() string {
	return m.ToOpenTSDBLine(false)
}

func (m *Metric) ToOpenTSDBLine(dropPut bool) string {
	res := []string{}
	if ! dropPut {
		res = append(res, "put")
	}
	res = append(res, m.Name, fmt.Sprintf("%d", m.Time.Unix()), fmt.Sprintf("%v", m.Value), m.GetDimensionString())
	return strings.Join(res, " ")
}

