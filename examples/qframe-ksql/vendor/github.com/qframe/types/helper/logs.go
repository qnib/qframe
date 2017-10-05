package qtypes_helper

import "github.com/qframe/types/constants"

func LogStrToInt(level string) int {
	def := 6
	switch level {
	case qtypes_constants.LOG_PANIC:
		return qtypes_constants.LOG_PANIC_INT
	case qtypes_constants.LOG_ERROR:
		return qtypes_constants.LOG_ERROR_INT
	case qtypes_constants.LOG_WARN:
		return qtypes_constants.LOG_WARN_INT
	case qtypes_constants.LOG_NOTICE:
		return qtypes_constants.LOG_NOTICE_INT
	case qtypes_constants.LOG_INFO:
		return qtypes_constants.LOG_INFO_INT
	case qtypes_constants.LOG_DEBUG:
		return qtypes_constants.LOG_DEBUG_INT
	case qtypes_constants.LOG_TRACE:
		return qtypes_constants.LOG_TRACE_INT
	default:
		return def
	}
}