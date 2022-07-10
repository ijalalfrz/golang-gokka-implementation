package wallet

import (
	"github.com/ijalalfrz/coinbit-test/model"
	"github.com/ijalalfrz/coinbit-test/pubsub"
	"github.com/lovoo/goka"
	"github.com/sirupsen/logrus"
)

// ProcessThresholdEventHandler is a concrete struct of wallet event handler.
type ProcessThresholdEventHandler struct {
	logger   *logrus.Logger
	usescase Usecase
}

// NewProcessThresholdEventHandler is a constructor.
func NewProcessThresholdEventHandler(logger *logrus.Logger, usecase Usecase) pubsub.GokaEventHandler {
	return &ProcessThresholdEventHandler{logger, usecase}
}

// Handle will process the message.
func (handler ProcessThresholdEventHandler) Handle(ctx goka.Context, message interface{}) {

	payload, ok := message.(*model.DepositWallet)
	if !ok {
		handler.logger.Error("Not a kafka message")
		return
	}
	result := handler.usescase.ProcessThreshold(ctx, payload)
	if result != nil {
		handler.logger.Info(result)
	}

	return
}
