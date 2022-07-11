package wallet_test

import (
	"testing"

	"github.com/ijalalfrz/coinbit-test/model"
	pubsubMock "github.com/ijalalfrz/coinbit-test/pubsub/mocks"
	"github.com/ijalalfrz/coinbit-test/response"
	"github.com/ijalalfrz/coinbit-test/wallet"
	"github.com/ijalalfrz/coinbit-test/wallet/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestOnDepositEventHandler_Error_When_CastMessage(t *testing.T) {
	usecase := mocks.Usecase{}
	context := pubsubMock.GokaContext{}
	handler := wallet.NewDepositWalletEventHandler(logrus.New(), &usecase)

	t.Run("Should error not a kafka message", func(t *testing.T) {
		handler.Handle(&context, nil)
	})
}

func TestOnDepositEventHandler_Success(t *testing.T) {
	usecase := mocks.Usecase{}
	context := pubsubMock.GokaContext{}
	handler := wallet.NewDepositWalletEventHandler(logrus.New(), &usecase)
	resp := response.NewSuccessResponse(nil, response.StatOK, "success")
	usecase.On("AddBalance", mock.Anything, mock.Anything).Return(resp)
	t.Run("Should error not a kafka message", func(t *testing.T) {
		payload := &model.DepositWallet{
			WalletId: "1",
			Amount:   1000,
		}
		handler.Handle(&context, payload)
	})
}
