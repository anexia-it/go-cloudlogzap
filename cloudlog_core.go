package cloudlogzap

import (
	"encoding/json"

	"github.com/anexia-it/go-cloudlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.Core = (*CloudLogCore)(nil)

// CloudlogClient interface allows you to pass your own implementation of a cloudlog client or mock clients
type CloudlogClient interface {
	PushEvent(interface{}) error
}

// CloudLogCore provides a custom zapcore.Core implementation for sending log messages to CloudLog
type CloudLogCore struct {
	client                CloudlogClient
	cloudLogClientOptions []cloudlog.Option
	cloudLogIndex         string
	parent                *zap.Logger

	zapcore.Core
}

type document struct {
	Message string                 `cloudlog:"message"`
	Level   string                 `cloudlog:"level"`
	Fields  map[string]interface{} `cloudlog:"fields"`
}

var convertFunc = func(entry zapcore.Entry, ff []zapcore.Field) interface{} {
	d := document{
		//Timestamp: time.Now().UTC().UnixNano() / int64(time.Millisecond),
		Message: entry.Message,
		Level:   entry.Level.String(),
	}

	config := zapcore.EncoderConfig{
		NameKey:        "module",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	enc := zapcore.NewJSONEncoder(config)
	buf, err := enc.EncodeEntry(entry, ff)
	if err != nil {
		return d
	}

	fields := make(map[string]interface{})
	if err = json.Unmarshal(buf.Bytes(), &fields); err != nil {
		return d
	}

	d.Fields = fields
	return d
}

// Check overrides the zapcore.Core Check method
func (cc *CloudLogCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(e, cc)
}

// Write overrides the zapcore.Core Write method
func (cc *CloudLogCore) Write(e zapcore.Entry, ff []zapcore.Field) (err error) {

	event := convertFunc(e, ff)
	err = cc.client.PushEvent(event)
	if err != nil {
		cc.parent.Debug("Write failed", zap.Error(err))
	}
	return
}

// NewCloudlogCore returns a new CloudLogCore or an error if no cloudlog.Client could be instantiated
func NewCloudlogCore(c zapcore.Core, index string, options []cloudlog.Option) (clc *CloudLogCore, err error) {
	var client *cloudlog.CloudLog
	client, err = cloudlog.NewCloudLog(index, options...)
	if err != nil {
		return
	}

	clc = &CloudLogCore{
		Core:                  c,
		cloudLogIndex:         index,
		cloudLogClientOptions: options,
		client:                client,
	}

	return
}
