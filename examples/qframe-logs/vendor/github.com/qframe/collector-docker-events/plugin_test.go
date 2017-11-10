package qcollector_docker_events

import (
	"testing"
	"github.com/go-test/deep"

	"github.com/zpatrick/go-config"
	"github.com/stretchr/testify/assert"
	"github.com/qframe/types/qchannel"
	"github.com/qframe/types/plugin"
)


func TestNew(t *testing.T) {
	qChan := qtypes_qchannel.NewQChan()
	kv := map[string]string{"log.level": "trace"}
	cfg := config.NewConfig([]config.Provider{config.NewStatic(kv)})
	exp := Plugin{
		Plugin: qtypes_plugin.NewNamedPlugin(qChan, cfg, pluginTyp, pluginPkg, "plugin", version),
	}
	got, err := New(qChan, cfg, "plugin")
	assert.NoError(t, err, "No error expected here")
	if diff := deep.Equal(exp, got); diff != nil {
		t.Error(diff)
	}
}