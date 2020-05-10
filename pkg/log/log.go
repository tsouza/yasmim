package log

type DelegateFactory func(category string) Delegate

type Delegate interface {
	Log(level Level, args ...interface{})
	Logf(level Level, format string, args ...interface{})
}

type Level int

const (
	LevelDebug  Level	= iota
	LevelInfo
	LevelWarn
)

func NewLogger(p Delegate) *Logger {
	return &Logger{p}
}

type Logger struct {
	impl Delegate
}

func (l *Logger) Debug(args ...interface{}) { l.impl.Log(LevelDebug, args...) }
func (l *Logger) Debugf(format string, args ...interface{}) { l.impl.Logf(LevelDebug, format, args...) }
func (l *Logger) Info(args ...interface{}) { l.impl.Log(LevelInfo, args...) }
func (l *Logger) Infof(format string, args ...interface{}) { l.impl.Logf(LevelInfo, format, args...) }
func (l *Logger) Warn(args ...interface{}) { l.impl.Log(LevelWarn, args...) }
func (l *Logger) Warnf(format string, args ...interface{}) { l.impl.Logf(LevelWarn, format, args...) }


func NoOpDelegateFactory(_ string) Delegate {
	return &noOpDelegate{}
}
type noOpDelegate struct {}
var _ Delegate = (*noOpDelegate)(nil)
func (d *noOpDelegate) Log(_ Level, _ ...interface{}) {}
func (d *noOpDelegate) Logf(_ Level, format string, _ ...interface{}) {}