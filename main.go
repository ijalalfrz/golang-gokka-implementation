package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.elastic.co/apm"

	"github.com/go-playground/validator/v10"
	"github.com/ijalalfrz/coinbit-test/pubsub"
	"github.com/ijalalfrz/coinbit-test/wallet"
	"github.com/lovoo/goka"

	"github.com/ijalalfrz/coinbit-test/middleware"

	"github.com/ijalalfrz/coinbit-test/response"
	"github.com/ijalalfrz/coinbit-test/server"

	gctx "github.com/gorilla/context"
	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"

	"github.com/ijalalfrz/coinbit-test/config"

	_ "github.com/joho/godotenv/autoload" // for development
	"go.elastic.co/apm/module/apmlogrus"
	_ "go.elastic.co/apm/module/apmsql/mysql"
)

var (
	tracer       *apm.Tracer
	cfg          *config.Config
	location     *time.Location
	tmc          *goka.TopicManagerConfig
	depositTopic string = "deposits"
	indexMessage string = "Application is running properly"
)

func init() {
	tracer = apm.DefaultTracer
	tmc = goka.NewTopicManagerConfig()
	tmc.Table.Replication = 1
	tmc.Stream.Replication = 1
	cfg = config.Load()
}

func main() {
	// init logger
	logger := logrus.New()
	logger.SetFormatter(cfg.Logger.Formatter)
	logger.SetReportCaller(true)
	logger.AddHook(&apmlogrus.Hook{
		LogLevels: logrus.AllLevels,
	})

	// init validator
	vld := validator.New()

	// init router object
	router := mux.NewRouter()
	router.HandleFunc("/wallet", index)

	// init topic event
	tm, err := goka.NewTopicManager(cfg.SaramaKafka.Addresses, goka.DefaultConfig(), tmc)
	if err != nil {
		logger.Fatalf("Error creating topic manager: %v", err)
	}
	defer tm.Close()
	err = tm.EnsureStreamExists(depositTopic, 1)
	if err != nil {
		logger.Fatal("Error creating kafka topic %s: %v", depositTopic, err)
	}
	// init codec for encode and decode
	depositWalletCodec := wallet.NewDepositCodec()
	walletCodec := wallet.NewWalletCodec()
	thresholdCodec := wallet.NewThresholdCodec()

	// init view table
	balanceVt := pubsub.NewGokaViewTableAdapter(logger, "balance", cfg.SaramaKafka.Addresses, walletCodec)
	thresholdVt := pubsub.NewGokaViewTableAdapter(logger, "aboveThreshold", cfg.SaramaKafka.Addresses, thresholdCodec)

	// init publisher
	depositTopicPublisher := pubsub.NewGokaProducerAdapter(logger, cfg.SaramaKafka.Addresses, depositTopic, depositWalletCodec)

	// init domain object
	walletUsecase := wallet.NewWalletUsecase(wallet.UsecaseProperty{
		ServiceName:           cfg.Application.Name,
		Logger:                logger,
		DepositTopicPublisher: depositTopicPublisher,
		RollingPeriod:         180,
		Threshold:             10000,
		BalanceViewTable:      balanceVt,
		ThresholdViewTable:    thresholdVt,
	})

	// init pub sub event
	depositWalletEventHandler := wallet.NewDepositWalletEventHandler(logger, walletUsecase)
	processThresholdEventHandler := wallet.NewProcessThresholdEventHandler(logger, walletUsecase)

	depositWalletBalanceGroup, err := pubsub.NewGokaConsumerGroupFullConfigAdapter(logger, cfg.SaramaKafka.Addresses,
		"balance", depositTopic, depositWalletEventHandler, tmc, depositWalletCodec, walletCodec)

	if err != nil {
		logger.Fatal(err)
	}

	processThresholdGroup, err := pubsub.NewGokaConsumerGroupFullConfigAdapter(logger, cfg.SaramaKafka.Addresses,
		"aboveThreshold", depositTopic, processThresholdEventHandler, tmc, depositWalletCodec, thresholdCodec)

	if err != nil {
		logger.Fatal(err)
	}

	// init http handler
	wallet.NewWalletHTTPHandler(logger, vld, router, walletUsecase)

	// middleware]
	httpHandler := gctx.ClearHandler(router)
	httpHandler = middleware.Recovery(logger, httpHandler)
	httpHandler = middleware.CORS(httpHandler)

	// initiate server
	srv := server.NewServer(logger, httpHandler, cfg.Application.Port)
	srv.Start()
	depositWalletBalanceGroup.Subscribe()
	processThresholdGroup.Subscribe()
	balanceVt.Open()
	thresholdVt.Open()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	<-sigterm

	// closing service for a gracefull shutdown.
	srv.Close()
	depositWalletBalanceGroup.Close()
	processThresholdGroup.Close()
	depositTopicPublisher.Close()
	balanceVt.Close()
	thresholdVt.Close()
}

func index(w http.ResponseWriter, r *http.Request) {
	resp := response.NewSuccessResponse(nil, response.StatOK, indexMessage)
	response.JSON(w, resp)
}
