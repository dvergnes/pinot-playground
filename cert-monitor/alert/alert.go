package alert

import (
	"bytes"
	"encoding/json"
)

// ObjectRef contains the information to locate the object in k8s, namely its name and namespace
type ObjectRef struct {
	// Name is the name of the k8s object
	Name string `json:"name"`
	// Namespace is the namespace of the k8s object
	Namespace string `json:"namespace"`
}

// Level defines the level of an alert, it is an enum of UNKNOWN, INFO, WARN, ERROR
type Level uint8

const (
	Unknown Level = iota
	Info
	Warn
	Error
)

// String returns a string representation of the Level value
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

// MarshalJSON marshals the enum as a quoted json string
func (l Level) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(l.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (l *Level) UnmarshalJSON(b []byte) error {
	var raw string
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}
	var val Level
	switch raw {
	case "INFO":
		val = Info
	case "WARN":
		val = Warn
	case "ERROR":
		val = Error
	default:
		val = Unknown
	}
	*l = val
	return nil
}

// Alert contains the information about the alert
type Alert struct {
	// Level defines the level of the alert
	Level Level `json:"level"`
	// Message describes the alert
	Message string `json:"message"`
	// ObjectRef defines the k8s object designated by the alert
	ObjectRef ObjectRef `json:"objectRef"`
	// Source defines the source of the alert
	Source string `json:"source"`
	// When defines when the alert has been created
	When int64 `json:"when"`
}

// Notifier is responsible to send an alert to an external system
type Notifier interface {
	// Send sends the alert to the external system
	Send(alert Alert) error
	// Close closes the notifier
	Close() error
}