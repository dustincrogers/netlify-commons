package service

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceConfig(t *testing.T) {
	cmd := new(cobra.Command)
	args := AddRootArgs(cmd)
	require.NotNil(t, args)

	f, err := ioutil.TempFile("", "svc-test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	envcontents :=
		`
LOG_LEVEL=debug
LOG_FILE=
LOG_FIELDS=abc:def,xyz:123
TRACING_ENABLED=false
SERVICE_FIELD=5
`
	require.NoError(t, ioutil.WriteFile(f.Name(), []byte(envcontents), 0644))
	args.envFile = f.Name()

	svcConfig := struct {
		ServiceField int `split_words:"true"`
	}{}

	log, err := args.InitService(t.Name(), "somegitsha", &svcConfig)
	require.NoError(t, err)
	require.NotNil(t, log)

	entry, ok := log.(*logrus.Entry)
	if assert.True(t, ok) {
		assert.Equal(t, logrus.DebugLevel, entry.Logger.Level)
		assert.Len(t, entry.Data, 3)
		assert.EqualValues(t, "def", entry.Data["abc"])
		assert.EqualValues(t, "123", entry.Data["xyz"])
		assert.EqualValues(t, "somegitsha", entry.Data["version"])
	}

	assert.Equal(t, opentracing.NoopTracer{}, opentracing.GlobalTracer())
}

func TestServiceConfigNothing(t *testing.T) {
	cmd := new(cobra.Command)
	args := AddRootArgs(cmd)
	require.NotNil(t, args)

	log, err := args.InitService("something", "", nil)
	require.NoError(t, err)
	require.NotNil(t, log)
	entry, ok := log.(*logrus.Entry)
	if assert.True(t, ok) {
		assert.Equal(t, logrus.InfoLevel, entry.Logger.Level)
		assert.Len(t, entry.Data, 1)
		assert.EqualValues(t, "", entry.Data["version"])
	}
}

func TestServiceConfigPrefix(t *testing.T) {
	cmd := new(cobra.Command)
	args := AddRootArgs(cmd)
	require.NotNil(t, args)

	f, err := ioutil.TempFile("", "svc-test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	envcontents :=
		`
TEST_LOG_LEVEL=debug
TEST_LOG_FILE=
TEST_LOG_FIELDS=abc:def,xyz:123
TEST_TRACING_ENABLED=true
TEST_TRACING_HOST=localhost
TEST_SERVICE_FIELD=5
`
	require.NoError(t, ioutil.WriteFile(f.Name(), []byte(envcontents), 0644))
	args.envFile = f.Name()
	args.prefix = "test"

	svcConfig := struct {
		ServiceField int `split_words:"true"`
	}{}

	log, err := args.InitService(t.Name(), "somegitsha", &svcConfig)
	require.NoError(t, err)
	require.NotNil(t, log)

	entry, ok := log.(*logrus.Entry)
	if assert.True(t, ok) {
		assert.Equal(t, logrus.DebugLevel, entry.Logger.Level)
		assert.Len(t, entry.Data, 3)
		assert.EqualValues(t, "def", entry.Data["abc"])
		assert.EqualValues(t, "123", entry.Data["xyz"])
		assert.EqualValues(t, "somegitsha", entry.Data["version"])
	}

	assert.NotEqual(t, opentracing.NoopTracer{}, opentracing.GlobalTracer())
}
