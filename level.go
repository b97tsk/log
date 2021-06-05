package log

// A Level is a logging level.
type Level int8

// LevelNone is the highest Level that doesn't allow to log anything.
const LevelNone Level = 127

// These are defined Levels, including LevelNone.
const (
	LevelError Level = 2 - iota
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

// Default logging level is LevelInfo, which must be zero.
const _, _ = uint8(LevelInfo), uint8(-LevelInfo)

// String returns the standard name for lv.
// String panics if lv is not a defined Level.
func (lv Level) String() string {
	switch lv {
	case LevelNone:
		return "NONE"
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelTrace:
		return "TRACE"
	}

	panic("unknown logging level")
}
