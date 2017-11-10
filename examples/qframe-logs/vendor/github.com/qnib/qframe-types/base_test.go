package qtypes

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestNewBase(t *testing.T) {
	before := time.Now()
	b := NewBase("src1")
	after := time.Now()
	assert.Equal(t, version, b.BaseVersion)
	assert.Equal(t, "src1", b.SourcePath[0])
	assert.True(t, before.UnixNano() < b.Time.UnixNano())
	assert.True(t, after.UnixNano() > b.Time.UnixNano())
}


func TestNewTimedBase(t *testing.T) {
	now := time.Now()
	b := NewTimedBase("src1", now)
	assert.Equal(t, now, b.Time)
}

func TestBase_GetTimeUnix(t *testing.T) {
	now := time.Now()
	b := NewTimedBase("src1", now)
	assert.Equal(t, now.Unix(), b.GetTimeUnix())
}

func TestBase_GetTimeUnixNano(t *testing.T) {
	now := time.Now()
	b := NewTimedBase( "src1", now)
	assert.Equal(t, now.UnixNano(), b.GetTimeUnixNano())
}


func TestBase_AppendSrc(t *testing.T) {
	b := NewBase("src1")
	b.AppendSource("src2")
	assert.Equal(t, "src1", b.SourcePath[0])
	assert.Equal(t, "src2", b.SourcePath[1])
}

func TestBase_IsLastSource(t *testing.T) {
	b := NewBase("src1")
	assert.True(t, b.IsLastSource("src1"), "Last source should be 'src1'")
	b.AppendSource("src2")
	assert.True(t, b.IsLastSource("src2"), "Last source should be 'src2'")
}

func TestBase_InputsMatch(t *testing.T) {
	b := NewBase("src1")
	assert.True(t, b.InputsMatch([]string{"src2", "src1"}), "Should match input list 'src2', 'src1'")
}

func TestSha1HashString(t *testing.T) {
	s := "sha1 this string"
	assert.Equal(t, "cf23df2207d99a74fbe169e3eba035e633b65d94", Sha1HashString(s))
}
