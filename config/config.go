package config

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.elastic.co/apm/module/apmmongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-redis/redis/v8"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

// Config is an app configuration.
type Config struct {
	Redis struct {
		Options *redis.Options
	}
	BasicAuth struct {
		Username string
		Password string
	}
	Application struct {
		Port string
		Name string
	}
	Crypto struct {
		Secret string
		IV     string
	}
	Logger struct {
		Formatter logrus.Formatter
	}
	Mariadb struct {
		Driver             string
		Host               string
		Port               string
		Username           string
		Password           string
		Database           string
		DSN                string
		MaxOpenConnections int
		MaxIdleConnections int
	}
	Mongodb struct {
		ClientOptions *options.ClientOptions
		Database      string
	}
	SaramaKafka struct {
		Addresses []string
		Config    *sarama.Config
	}
	ThirdPartyAggregation struct {
		Host string
	}
}

// Load will load the configuration.
func Load() *Config {
	cfg := new(Config)
	cfg.mongodb()
	cfg.redis()
	cfg.basicAuth()
	cfg.sarama()
	cfg.crypto()
	cfg.logFormatter()
	cfg.app()
	cfg.mariadb()
	cfg.thirdPartyAggregation()
	return cfg
}

func (cfg *Config) thirdPartyAggregation() {
	host := os.Getenv("THIRD_PARTY_MANAGEMENT_AGGREGATION_SERVICE")

	cfg.ThirdPartyAggregation.Host = host
}

func (cfg *Config) mongodb() {
	appName := os.Getenv("APP_NAME")
	uri := os.Getenv("MONGODB_URL")
	db := os.Getenv("MONGODB_DATABASE")
	minPoolSize, _ := strconv.ParseUint(os.Getenv("MONGODB_MIN_POOL_SIZE"), 10, 64)
	maxPoolSize, _ := strconv.ParseUint(os.Getenv("MONGODB_MAX_POOL_SIZE"), 10, 64)
	maxConnIdleTime, _ := strconv.ParseInt(os.Getenv("MONGODB_MAX_IDLE_CONNECTION_TIME_MS"), 10, 64)

	opts := options.Client().
		ApplyURI(uri).
		SetAppName(appName).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize).
		SetMaxConnIdleTime(time.Millisecond * time.Duration(maxConnIdleTime)).
		SetMonitor(apmmongo.CommandMonitor())

	cfg.Mongodb.ClientOptions = opts
	cfg.Mongodb.Database = db
}
func (cfg *Config) redis() {
	host := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.ParseInt(os.Getenv("REDIS_DATABASE"), 10, 64)

	options := &redis.Options{
		Addr:     host,
		Password: password,
		DB:       int(db),
	}

	cfg.Redis.Options = options
}

func (cfg *Config) basicAuth() {
	username := os.Getenv("BASIC_AUTH_USERNAME")
	password := os.Getenv("BASIC_AUTH_PASSWORD")

	cfg.BasicAuth.Username = username
	cfg.BasicAuth.Password = password
}

func (cfg *Config) sarama() {
	brokers := os.Getenv("KAFKA_BROKERS")
	sslEnable, _ := strconv.ParseBool(os.Getenv("KAFKA_SSL_ENABLE"))
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")

	sc := sarama.NewConfig()
	sc.Version = sarama.V2_1_0_0
	if username != "" {
		sc.Net.SASL.User = username
		sc.Net.SASL.Password = password
		sc.Net.SASL.Enable = true
	}
	sc.Net.TLS.Enable = sslEnable

	// consumer config
	sc.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	sc.Consumer.Offsets.Initial = sarama.OffsetOldest

	// producer config
	sc.Producer.Retry.Backoff = time.Millisecond * 500

	cfg.SaramaKafka.Addresses = strings.Split(brokers, ",")
	cfg.SaramaKafka.Config = sc
}

func (cfg *Config) crypto() {
	secret := os.Getenv("AES_SECRET")
	iv := os.Getenv("AES_IV")

	cfg.Crypto.IV = iv
	cfg.Crypto.Secret = secret
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

func (cfg *Config) mariadb() {
	host := os.Getenv("MARIADB_HOST")
	port := os.Getenv("MARIADB_PORT")
	username := os.Getenv("MARIADB_USERNAME")
	password := os.Getenv("MARIADB_PASSWORD")
	database := os.Getenv("MARIADB_DATABASE")
	maxOpenConnections, _ := strconv.ParseInt(os.Getenv("MARIADB_MAX_OPEN_CONNECTIONS"), 10, 64)
	maxIdleConnections, _ := strconv.ParseInt(os.Getenv("MARIADB_MAX_IDLE_CONNECTIONS"), 10, 64)

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	connVal := url.Values{}
	connVal.Add("parseTime", "1")
	connVal.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", dbConnectionString, connVal.Encode())

	cfg.Mariadb.Driver = "mysql"
	cfg.Mariadb.Host = host
	cfg.Mariadb.Port = port
	cfg.Mariadb.Username = username
	cfg.Mariadb.Password = password
	cfg.Mariadb.Database = database
	cfg.Mariadb.DSN = dsn
	cfg.Mariadb.MaxOpenConnections = int(maxOpenConnections)
	cfg.Mariadb.MaxIdleConnections = int(maxIdleConnections)
}
