package wallet

import (
	"github.com/ijalalfrz/coinbit-test/model"
	"github.com/ijalalfrz/coinbit-test/pubsub"
	"github.com/lovoo/goka"
	"github.com/sirupsen/logrus"
)

// DepositWalletEventHandler is a concrete struct of wallet event handler.
type DepositWalletEventHandler struct {
	logger   *logrus.Logger
	usescase Usecase
}

// NewDepositWalletEventHandler is a constructor.
func NewDepositWalletEventHandler(logger *logrus.Logger, usecase Usecase) pubsub.GokaEventHandler {
	return &DepositWalletEventHandler{logger, usecase}
}

// Handle will process the message.
func (handler DepositWalletEventHandler) Handle(ctx goka.Context, message interface{}) {

	payload, ok := message.(*model.DepositWallet)
	if !ok {
		handler.logger.Error("Not a kafka message")
		return
	}

	result := handler.usescase.AddBalance(ctx, payload)
	if result != nil {
		handler.logger.Info(result)
	}

	return
}
