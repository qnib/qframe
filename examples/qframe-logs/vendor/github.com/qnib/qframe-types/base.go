package qtypes

import (
	"time"
	"crypto/sha1"
	"fmt"
)

const (
	version = "0.5.11"
)

type Base struct {
	BaseVersion string
	ID				string
	Time			time.Time
	SourceID		int
	SourcePath		[]string
	SourceSuccess 	bool
	Data 			map[string]string // Additional Data

}

func NewBase(src string) Base {
	return NewTimedBase(src, time.Now())
}

func NewTimedBase(src string, t time.Time) Base {
	b := Base {
		BaseVersion: version,
		ID: "",
		Time: t,
		SourceID: 0,
		SourcePath: []string{src},
		SourceSuccess: true,
		Data: map[string]string{},
	}
	return b
}

// GenDefaultID uses "<source>-<time.UnixNano()>" and does a sha1 hash.
func (b *Base) GenDefaultID() string {
	s := fmt.Sprintf("%s-%d", b.GetLastSource(), b.Time.UnixNano())
	return Sha1HashString(s)
}

func (b *Base) GetMessageDigest() string {
	return b.ID[:13]
}

func (base *Base) NewExtMetric(src, name string, metricTyp string, val float64, dimensions map[string]string, t time.Time, buffered bool) Metric {
	m := Metric{
		Base: 		*base,
		Name:       sanitizeString(name),
		MetricType: metricTyp,
		Value:      val,
		Dimensions: dimensions,
		Buffered:   buffered,
	}
	m.AppendSource(src)
	return m
}

func (b *Base) GetTimeRFC() string {
	return b.Time.Format("2006-01-02T15:04:05.999999-07:00")
}

func (b *Base) GetTimeUnix() int64 {
	return b.Time.Unix()
}

func (b *Base) GetTimeUnixNano() int64 {
	return b.Time.UnixNano()
}

func (b *Base) AppendSource(src string) {
	b.SourcePath = append(b.SourcePath, src)
}

func (b *Base) GetLastSource() string {
	return b.SourcePath[len(b.SourcePath)-1]
}

func (b *Base) IsLastSource(src string) bool {
	return b.SourcePath[len(b.SourcePath)-1] == src
}

func (b *Base) InputsMatch(inputs []string) bool {
	for _, inp := range inputs {
		if b.IsLastSource(inp) {
			return true
		}

	}
	return false
}

func Sha1HashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}