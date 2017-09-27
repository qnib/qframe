package qtypes_syslog

import (
	"text/template"
	"bytes"
)

const (
	KEY_PRI = "syslog5424_pri"
	KEY_VER = "syslog5424_ver"
	KEY_TIME = "syslog5424_ts"
	KEY_HOST = "syslog5424_host"
	KEY_APP = "syslog5424_app"
	KEY_PROC = "syslog5424_proc"
	KEY_MSGID = "syslog5424_msgid"
	KEY_SD = "syslog5424_sd"
	KEY_MSG = "syslog5424_msg"
	TEMPLATE = `<{{.Pri}}>{{.UserVer}} {{.Time}} {{.Host}} {{if eq .App ""}}-{{else}}{{.App}}{{end}} {{if eq .Proc ""}}-{{else}}{{.Proc}}{{end}} {{if eq .MsgID ""}}-{{else}}{{.MsgID}}{{end}} {{if eq .Structured ""}}{{else}} {{.Structured}} {{end}}{{if .IsCee}}@cee:{{end}}{{.Message}}`
)

type Syslog struct {
	Pri string
	UserVer string
	Time string
	Host, App, Proc, MsgID string
	Structured string
	Message string
	IsCee bool
}

func NewSyslogFromKV(kv map[string]string) (sl Syslog, err error) {
	sl.Pri = kv[KEY_PRI]
	sl.UserVer = kv[KEY_VER]
	sl.Time = kv[KEY_TIME]
	sl.Host = kv[KEY_HOST]
	sl.App = kv[KEY_APP]
	sl.Proc = kv[KEY_PROC]
	sl.MsgID = kv[KEY_MSGID]
	sl.Structured = kv[KEY_SD]
	sl.Message = kv[KEY_MSG]
	sl.IsCee = false
	return
}

func (sl *Syslog) EnableCEE() {
	sl.IsCee = true
}


func (sl *Syslog) SetMessage(msg string) {
	sl.Message = msg
}

func (sl *Syslog) ToRFC5424() (str string, err error) {
	var tpl bytes.Buffer
	tmpl, err := template.New("test").Parse(TEMPLATE)
	err = tmpl.Execute(&tpl, sl)
	return tpl.String(), err
}
