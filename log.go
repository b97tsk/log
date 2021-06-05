// Package log is a library for logging messages.
package log

import (
	"bytes"
	"io"
	"log"
	"sync"
)

// A Writer is an io.Writer and also provides Writable method to check if it's
// writable at certain Level.
type Writer interface {
	io.Writer
	Writable(Level) bool
}

// A Logger represents an active logging object that generates lines of output
// to a Writer. Each logging operation makes a single call to the Writer's
// Write method. Unless the associated Writer is concurrency-safe, a Logger
// cannot be used simultaneously from multiple goroutines.
type Logger struct {
	Writer

	e *log.Logger
	w *log.Logger
	i *log.Logger
	d *log.Logger
	t *log.Logger
}

const levelPlaceHolder = "[#####]"

// New creates a new Logger. The out variable sets the destination to which log
// data will be written. The prefix appears at the beginning of each generated
// log line, or after the log header if the Lmsgprefix flag is provided. The
// flag argument defines the logging properties.
func New(out Writer, prefix string, flags int) *Logger {
	prefix = levelPlaceHolder + " " + prefix

	return &Logger{
		Writer: out,

		e: log.New(&writer{out, LevelError}, prefix, flags),
		w: log.New(&writer{out, LevelWarn}, prefix, flags),
		i: log.New(&writer{out, LevelInfo}, prefix, flags),
		d: log.New(&writer{out, LevelDebug}, prefix, flags),
		t: log.New(&writer{out, LevelTrace}, prefix, flags),
	}
}

var noneLogger = log.New(io.Discard, "", 0)

// Get returns a log.Logger at Level lv.
// Get panics if lv is not a defined Level.
func (l *Logger) Get(lv Level) *log.Logger {
	switch lv {
	case LevelNone:
		return noneLogger
	case LevelError:
		return l.e
	case LevelWarn:
		return l.w
	case LevelInfo:
		return l.i
	case LevelDebug:
		return l.d
	case LevelTrace:
		return l.t
	}

	panic("unknown logging level")
}

// Error logs a message at LevelError.
func (l *Logger) Error(v ...interface{}) {
	if l.ErrorWritable() {
		l.e.Print(v...)
	}
}

// Errorf logs a message at LevelError.
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.ErrorWritable() {
		l.e.Printf(format, v...)
	}
}

// Errorln logs a message at LevelError.
func (l *Logger) Errorln(v ...interface{}) {
	if l.ErrorWritable() {
		l.e.Println(v...)
	}
}

// ErrorWritable reports whether l can write messages at LevelError.
func (l *Logger) ErrorWritable() bool {
	return l.Writable(LevelError)
}

// Warn logs a message at LevelWarn.
func (l *Logger) Warn(v ...interface{}) {
	if l.WarnWritable() {
		l.w.Print(v...)
	}
}

// Warnf logs a message at LevelWarn.
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.WarnWritable() {
		l.w.Printf(format, v...)
	}
}

// Warnln logs a message at LevelWarn.
func (l *Logger) Warnln(v ...interface{}) {
	if l.WarnWritable() {
		l.w.Println(v...)
	}
}

// WarnWritable reports whether l can write messages at LevelWarn.
func (l *Logger) WarnWritable() bool {
	return l.Writable(LevelWarn)
}

// Info logs a message at LevelInfo.
func (l *Logger) Info(v ...interface{}) {
	if l.InfoWritable() {
		l.i.Print(v...)
	}
}

// Infof logs a message at LevelInfo.
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.InfoWritable() {
		l.i.Printf(format, v...)
	}
}

// Infoln logs a message at LevelInfo.
func (l *Logger) Infoln(v ...interface{}) {
	if l.InfoWritable() {
		l.i.Println(v...)
	}
}

// InfoWritable reports whether l can write messages at LevelInfo.
func (l *Logger) InfoWritable() bool {
	return l.Writable(LevelInfo)
}

// Debug logs a message at LevelDebug.
func (l *Logger) Debug(v ...interface{}) {
	if l.DebugWritable() {
		l.d.Print(v...)
	}
}

// Debugf logs a message at LevelDebug.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.DebugWritable() {
		l.d.Printf(format, v...)
	}
}

// Debugln logs a message at LevelDebug.
func (l *Logger) Debugln(v ...interface{}) {
	if l.DebugWritable() {
		l.d.Println(v...)
	}
}

// DebugWritable reports whether l can write messages at LevelDebug.
func (l *Logger) DebugWritable() bool {
	return l.Writable(LevelDebug)
}

// Trace logs a message at LevelTrace.
func (l *Logger) Trace(v ...interface{}) {
	if l.TraceWritable() {
		l.t.Print(v...)
	}
}

// Tracef logs a message at LevelTrace.
func (l *Logger) Tracef(format string, v ...interface{}) {
	if l.TraceWritable() {
		l.t.Printf(format, v...)
	}
}

// Traceln logs a message at LevelTrace.
func (l *Logger) Traceln(v ...interface{}) {
	if l.TraceWritable() {
		l.t.Println(v...)
	}
}

// TraceWritable reports whether l can write messages at LevelTrace.
func (l *Logger) TraceWritable() bool {
	return l.Writable(LevelTrace)
}

type writer struct {
	w  Writer
	lv Level
}

func (w *writer) Write(p []byte) (n int, err error) {
	if !w.w.Writable(w.lv) {
		return len(p), nil
	}

	if i := bytes.Index(p, []byte(levelPlaceHolder)); i >= 0 {
		b := pool.Get().(*buffer)
		defer pool.Put(b)

		s := (*b)[:0]
		s = append(s, p[:i]...)
		s = append(s, '[')
		s = append(s, w.lv.String()...)
		s = append(s, ']')
		s = append(s, p[i+len(levelPlaceHolder):]...)
		p = s
	}

	return w.w.Write(p)
}

type buffer [1024]byte

var pool = sync.Pool{
	New: func() interface{} { return new(buffer) },
}
