package config_test

import (
	"os"
	"testing"

	"github.com/ijalalfrz/coinbit-test/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	os.Setenv("KAFKA_USERNAME", "test_username")
	cfg := config.Load()

	assert.NotNil(t, cfg)

	t.Run("when config is being used for logger", func(t *testing.T) {
		logger := logrus.New()
		logger.SetFormatter(cfg.Logger.Formatter)
		logger.SetReportCaller(true)

		logger.Info("called")
	})

}
