package alert

import "go.uber.org/zap"

// NewLogNotifier returns a Notifier that generate logs using the passed logger
func NewLogNotifier(logger *zap.SugaredLogger) Notifier {
	return &logNotifier{
		logger: logger,
	}
}

// logNotifier generates logs when asked to send an alert
type logNotifier struct {
	logger *zap.SugaredLogger
}

// Send implements Notifier contract
func (l *logNotifier) Send(alert Alert) error {
	switch alert.Level {
	case Info:
		l.logger.Infow("processing notification",
			"message", alert.Message,
			"object", alert.ObjectRef,
			"when", alert.When,
			"source", alert.Source)
	case Warn:
		l.logger.Warnw("processing notification",
			"message", alert.Message,
			"object", alert.ObjectRef,
			"when", alert.When,
			"source", alert.Source)
	case Error:
		l.logger.Errorw("processing notification",
			"message", alert.Message,
			"object", alert.ObjectRef,
			"when", alert.When,
			"source", alert.Source)
	default:
		l.logger.Warnw("processing notification with unexpected level",
			"message", alert.Message,
			"object", alert.ObjectRef,
			"when", alert.When,
			"source", alert.Source,
			"level", alert.Level)
	}
	return nil
}
