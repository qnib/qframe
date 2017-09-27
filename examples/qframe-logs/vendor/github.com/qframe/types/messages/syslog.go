package qtypes_messages

import "github.com/qframe/types/syslog"

type SyslogMessage struct {
	Base
	Syslog qtypes_syslog.Syslog
}

func NewSyslogMessage(b Base, sl qtypes_syslog.Syslog) SyslogMessage {
	return SyslogMessage{
		Base: b,
		Syslog: sl,
	}
}

func (sm *SyslogMessage) ToRFC5424() (str string, err error) {
	str, err = sm.Syslog.ToRFC5424()
	return
}
