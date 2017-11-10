package qcache_inventory

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/zpatrick/go-config"
	"github.com/qframe/types/qchannel"
)

func TestNew(t *testing.T) {
	cfg := config.NewConfig([]config.Provider{})
	qChan := qtypes_qchannel.NewCfgQChan(cfg)
	_, err := New(qChan, cfg, "test")
	assert.NoError(t, err, "should not cause trouble")
}
