package alert

type ObjectRef struct {
	Namespace string `json:"namespace"`
	Name string `json:"name"`
}

type Level uint8

const (
	Unknown Level = iota
	Info
	Warn
	Error
)


func (l Level) String() string {
	switch l {
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Alert struct {
	Level Level `json:"level"`
	ObjectRef ObjectRef `json:"objectRef"`
	Message string `json:"message"`
	When int64 `json:"when"`
	Source string `json:"source"`
}

type Notifier interface {
	Send(alert Alert) error
}