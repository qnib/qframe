package qtypes_syslog

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestNewSyslogFromKV(t *testing.T) {
	kv := map[string]string{
		"syslog5424_pri": "14",
		"syslog5424_ver": "1",
		"syslog5424_ts": "2017-09-14T15:02:07.011Z",
		"syslog5424_host": "e4d0e5567436",
		"syslog5424_app": "-",
		"syslog5424_proc": "-",
		"syslog5424_msgid": "event",
		"syslog5424_sd": "",
		"syslog5424_msg": `Hello World`,
	}
	sl, err := NewSyslogFromKV(kv)
	assert.NoError(t, err, "Should assemble a fine syslog struct")
	msg, err := sl.ToRFC5424()
	exp := `<14>1 2017-09-14T15:02:07.011Z e4d0e5567436 - - event Hello World`
	assert.Equal(t, exp, msg)
	kv["syslog5424_msg"] = `{"message": "Hello World"}`
	sl, err = NewSyslogFromKV(kv)
	sl.EnableCEE()
	msg, err = sl.ToRFC5424()
	exp = `<14>1 2017-09-14T15:02:07.011Z e4d0e5567436 - - event @cee:{"message": "Hello World"}`
	assert.Equal(t, exp, msg)

}
