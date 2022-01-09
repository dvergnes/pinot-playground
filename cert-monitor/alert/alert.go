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


func (l Level) MarshalJSON() ([]byte, error) {
	switch l {
	case Info:
		return []byte("INFO"), nil
	case Warn:
		return []byte("WARN"), nil
	case Error:
		return []byte("ERROR"), nil
	default:
		return []byte("UNKNOWN"), nil
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