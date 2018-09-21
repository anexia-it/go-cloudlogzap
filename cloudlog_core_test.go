package cloudlogzap

import (
	"github.com/anexia-it/go-cloudlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

type MockCloudlogClient struct {
	events []interface{}
}

func (client *MockCloudlogClient) PushEvent(e interface{}) error {
	client.events = append(client.events, e)
	return nil
}

func TestNewCloudlogCore(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		core, err := NewCloudlogCore(zapcore.NewNopCore(), "testindex", nil)
		require.NoError(t, err)
		require.NotNil(t, core)
		assert.EqualValues(t, "testindex", core.cloudLogIndex)
		assert.NotNil(t, core.client)
	})

	t.Run("Nil", func(t *testing.T) {
		core, err := NewCloudlogCore(zapcore.NewNopCore(), "", nil)
		require.Error(t, err)
		assert.Nil(t, core)
	})
}

func TestNewCloudlogCore_WithOption(t *testing.T) {
	expected := []cloudlog.Option{cloudlog.OptionBrokers(cloudlog.DefaultBrokerAddresses...)}
	core, err := NewCloudlogCore(zapcore.NewNopCore(), "testindex", expected)
	require.NoError(t, err)
	require.NotNil(t, core)
	require.EqualValues(t, expected, core.cloudLogClientOptions)
}

func TestCloudLogCore_Check(t *testing.T) {
	core, err := NewCloudlogCore(zapcore.NewNopCore(), "testindex", nil)
	require.NoError(t, err)
	entry := zapcore.Entry{
		Time:       time.Now(),
		Level:      zapcore.InfoLevel,
		Stack:      "",
		Caller:     zapcore.EntryCaller{},
		LoggerName: "test",
		Message:    "test message",
	}

	checkedEntry := &zapcore.CheckedEntry{
		ErrorOutput: zapcore.NewMultiWriteSyncer(),
	}
	result := core.Check(entry, checkedEntry)
	require.NotNil(t, result)
}

func TestCloudLogCore_Write(t *testing.T) {
	core, err := NewCloudlogCore(zapcore.NewNopCore(), "testindex", nil)
	require.NoError(t, err)
	client := &MockCloudlogClient{}
	core.client = client
	entry := zapcore.Entry{
		Time:       time.Now(),
		Level:      zapcore.InfoLevel,
		Stack:      "",
		Caller:     zapcore.EntryCaller{},
		LoggerName: "test",
		Message:    "test message",
	}
	ff := make([]zapcore.Field, 0)
	err = core.Write(entry, ff)
	require.NoError(t, err)
	assert.Len(t, client.events, 1)

}

func TestCloudLogCore_ConverterFunc(t *testing.T) {
	entry := zapcore.Entry{
		Level:      zapcore.InfoLevel,
		LoggerName: "test",
		Message:    "test message",
	}
	ff := []zapcore.Field{
		{
			Key:    "key",
			String: "value",
			Type:   zapcore.StringType,
		},
		{
			Key:     "key2",
			Integer: 42,
			Type:    zapcore.Int64Type,
		},
	}
	result := convertFunc(entry, ff)
	d, ok := result.(document)
	require.True(t, ok)
	assert.EqualValues(t, "test message", d.Message)
	assert.EqualValues(t, "info", d.Level)
	for _, expected := range ff {
		_, ok := d.Fields[expected.Key]
		assert.True(t, ok)
	}
	assert.EqualValues(t, d.Fields["module"], entry.LoggerName)
}
