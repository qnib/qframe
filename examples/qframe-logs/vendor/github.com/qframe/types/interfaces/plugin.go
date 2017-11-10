package qtypes_interfaces


// QPlugin
type QPlugin interface {
	// Starts consumes bcast channels, config and attributes of a Plugin and runs a loop
	GetInfo() (typ,pkg,name string)
	CfgStringOr(path, alt string) string
	GetLogOnlyPlugs() []string
}
