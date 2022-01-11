// Copyright (c) 2022 Denis Vergnes
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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

// Close implements Notifier contract
func (l *logNotifier) Close() error {
	return nil
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
