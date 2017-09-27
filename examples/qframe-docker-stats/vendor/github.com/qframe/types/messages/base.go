package qtypes_messages

import (
	"time"
	"crypto/sha1"
	"fmt"
	"github.com/qframe/types/plugin"
)

const (
	version = "0.1.3"
)

type Base struct {
	BaseVersion string
	ID				string
	Time			time.Time
	SourceID		int
	SourcePath		[]string
	SourceSuccess 	bool
	Tags 			map[string]string // Additional KV
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
		Tags: map[string]string{},
	}
	return b
}

func NewBaseFromBase(src string, b Base) Base {
	return Base {
		BaseVersion: b.BaseVersion,
		ID: b.ID,
		Time: b.Time,
		SourceID: b.SourceID,
		SourcePath: append(b.SourcePath, src),
		SourceSuccess: b.SourceSuccess,
		Tags: b.Tags,
	}
}


func (b *Base) ToJSON() map[string]interface{} {
	res := map[string]interface{}{
		"base_version": b.BaseVersion,
		"id": b.ID,
		"time": b.Time.String(),
		"time_unix_nano": b.Time.UnixNano(),
	}
	res["source_id"] = b.SourceID
	res["source_path"] = b.SourcePath
	res["source_success"] = b.SourceSuccess
	res["tags"] = b.Tags
	return res
}

// GenDefaultID uses "<source>-<time.UnixNano()>" and does a sha1 hash.
func (b *Base) GenDefaultID() string {
	s := fmt.Sprintf("%s-%d", b.GetLastSource(), b.Time.UnixNano())
	return Sha1HashString(s)
}

func (b *Base) GetMessageDigest() string {
	return b.ID[:13]
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


func (b *Base) StopProcessing(p *qtypes_plugin.Plugin, allowEmptyInput bool) bool {
	if b.SourceID != 0 && p.MyID == b.SourceID {
		msg := fmt.Sprintf("Msg came from the same GID (My:%d == %d:SourceID)", p.MyID, b.SourceID)
		p.Log("debug", msg)
		return true
	}
	// TODO: Most likely invoked often, so check if performant enough
	inputs := p.GetInputs()
	if ! allowEmptyInput && len(inputs) == 0 {
		format := "Plugin '%s' does not allow empty imputs, please set '%s.%s.inputs'"
		msg := fmt.Sprintf(format, p.Name, p.Typ, p.Name)
		p.Log("error", msg)
		return true
	}
	srcSuccess := p.CfgBoolOr("source-success", true)
	if ! b.InputsMatch(inputs) {
		p.Log("debug", fmt.Sprintf("InputsMatch(%v) != %s", inputs, b.GetLastSource()))
		return true
	}
	if b.SourceSuccess != srcSuccess {
		msg := fmt.Sprintf("qm.SourceSuccess (%v) != (%v) srcSuccess", b.SourceSuccess, srcSuccess)
		p.Log("debug", msg)
		return true
	}
	return false
}
