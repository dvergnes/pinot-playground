package alert

import "go.uber.org/zap"

func NewLogNotifier(logger *zap.SugaredLogger) Notifier {
	return &LogNotifer{
		logger: logger,
	}
}

type LogNotifer struct {
	logger *zap.SugaredLogger
}

func (l *LogNotifer) Send(alert Alert) error {
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
