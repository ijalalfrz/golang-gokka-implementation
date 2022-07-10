package wallet

import (
	"github.com/ijalalfrz/coinbit-test/pubsub"
	"github.com/sirupsen/logrus"
)

type UsecaseProperty struct {
	ServiceName           string
	Logger                *logrus.Logger
	DepositTopicPublisher pubsub.Publisher
	RollingPeriod         int
	Threshold             int64
	BalanceViewTable      pubsub.ViewTable
	ThresholdViewTable    pubsub.ViewTable
}
