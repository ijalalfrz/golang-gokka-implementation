package config

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

// Config is an app configuration.
type Config struct {
	Application struct {
		Port string
		Name string
	}
	Logger struct {
		Formatter logrus.Formatter
	}
	SaramaKafka struct {
		Addresses []string
		Config    *sarama.Config
	}
	Wallet struct {
		Threshold     int64
		RollingPeriod int
	}
}

// Load will load the configuration.
func Load() *Config {
	cfg := new(Config)
	cfg.sarama()
	cfg.logFormatter()
	cfg.wallet()
	cfg.app()
	return cfg
}

func (cfg *Config) sarama() {
	brokers := os.Getenv("KAFKA_BROKERS")
	// sslEnable, _ := strconv.ParseBool(os.Getenv("KAFKA_SSL_ENABLE"))
	// username := os.Getenv("KAFKA_USERNAME")
	// password := os.Getenv("KAFKA_PASSWORD")

	// sc := sarama.NewConfig()
	// sc.Version = sarama.V2_1_0_0
	// if username != "" {
	// 	sc.Net.SASL.User = username
	// 	sc.Net.SASL.Password = password
	// 	sc.Net.SASL.Enable = true
	// }
	// sc.Net.TLS.Enable = sslEnable

	// // consumer config
	// sc.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	// sc.Consumer.Offsets.Initial = sarama.OffsetOldest

	// // producer config
	// sc.Producer.Retry.Backoff = time.Millisecond * 500

	cfg.SaramaKafka.Addresses = strings.Split(brokers, ",")
}

func (cfg *Config) logFormatter() {
	formatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	}

	cfg.Logger.Formatter = formatter
}

func (cfg *Config) app() {
	appName := os.Getenv("APP_NAME")
	port := os.Getenv("PORT")

	cfg.Application.Port = port
	cfg.Application.Name = appName
}

func (cfg *Config) wallet() {
	threshold, _ := strconv.ParseInt(os.Getenv("THRESHOLD"), 10, 64)
	rollingPeriod, _ := strconv.Atoi(os.Getenv("ROLLING_PERIOD"))

	cfg.Wallet.RollingPeriod = rollingPeriod
	cfg.Wallet.Threshold = threshold
}
