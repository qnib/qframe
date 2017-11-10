package qtypes

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/zpatrick/go-config"
)

func TestPlugin_CfgString(t *testing.T) {
	cfgMap := map[string]string{
		"testTyp.testName.database": "qframe",
	}
	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)
	p := NewNamedPlugin(NewQChan(), cfg, "testTyp", "testPkg", "testName", "0.0.0")
	got, err := p.CfgString("database")
	assert.NoError(t, err)
	assert.Equal(t, "qframe", got)
	got = p.CfgStringOr("database", "alt")
	assert.Equal(t, "qframe", got)
	got = p.CfgStringOr("nil", "alt")
	assert.Equal(t, "alt", got)
}

func TestPlugin_CfgBool(t *testing.T) {
	cfgMap := map[string]string{
		"testTyp.testName.key": "true",
		"testTyp.testName.key2": "nil",
	}
	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)
	p := NewNamedPlugin(NewQChan(), cfg, "testTyp", "testPkg", "testName", "0.0.0")
	got, err := p.CfgBool("key")
	assert.NoError(t, err)
	assert.True(t, got)
	got = p.CfgBoolOr("key", false)
	assert.True(t, got)
	got = p.CfgBoolOr("nil", false)
	assert.False(t, got)
	_, err = p.CfgBool("key2")
	assert.Error(t, err, "Should read 'nil'")

}

func TestPlugin_CfgInt(t *testing.T) {
	cfgMap := map[string]string{
		"testTyp.testName.key": "1",
	}
	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)
	p := NewNamedPlugin(NewQChan(), cfg, "testTyp", "testPkg", "testName", "0.0.0")
	got, err := p.CfgInt("key")
	assert.NoError(t, err)
	assert.Equal(t, 1, got)
	got = p.CfgIntOr("key", 2)
	assert.Equal(t, 1, got)
	got = p.CfgIntOr("nil", 2)
	assert.Equal(t, 2, got)
}

func TestPlugin_GetInputs(t *testing.T) {
	cfg := config.NewConfig([]config.Provider{})
	p := NewNamedPlugin(NewQChan(), cfg, "testTyp", "testPkg", "testName", "0.0.0")
	assert.Equal(t, 0, len(p.GetInputs()))
	cfg = config.NewConfig([]config.Provider{config.NewStatic(map[string]string{
		"testTyp.testName.inputs": "test",
	})})
	p = NewNamedPlugin(NewQChan(), cfg, "testTyp", "testPkg", "testName", "0.0.0")
	assert.Equal(t, 1, len(p.GetInputs()))
}