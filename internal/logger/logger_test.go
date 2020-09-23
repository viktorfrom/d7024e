package logger

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateNew(t *testing.T) {
	filename := "hello.log"
	brokenFilename := ""
	logger1 := New(log.DebugLevel, nil, false)
	logger2 := New(log.ErrorLevel, &filename, true)
	logger3 := New(log.InfoLevel, &brokenFilename, true)

	assert.NotNil(t, logger1)
	assert.NotNil(t, logger2)
	assert.NotNil(t, logger3)

	os.Remove(filename)
}

func TestLogInfo(t *testing.T) {
	filename := "info.log"
	logger := New(log.InfoLevel, &filename, false)
	logger.Info("info text here")

	file, err := os.Stat(filename)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, file.Size())

	os.Remove(filename)
}

func TestLogWarning(t *testing.T) {
	filename := "warning.log"
	logger := New(log.WarnLevel, &filename, false)
	logger.Warning("warning text here")

	file, err := os.Stat(filename)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, file.Size())

	os.Remove(filename)
}

func TestLogError(t *testing.T) {
	filename := "error.log"
	logger := New(log.ErrorLevel, &filename, false)
	logger.Error("error text here")

	file, err := os.Stat(filename)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, file.Size())

	os.Remove(filename)
}
