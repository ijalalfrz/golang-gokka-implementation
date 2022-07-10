package wallet_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ijalalfrz/coinbit-test/entity"
	"github.com/ijalalfrz/coinbit-test/exception"
	"github.com/ijalalfrz/coinbit-test/model"
	pubsubMock "github.com/ijalalfrz/coinbit-test/pubsub/mocks"
	"github.com/ijalalfrz/coinbit-test/response"
	"github.com/ijalalfrz/coinbit-test/wallet"
	"github.com/ijalalfrz/coinbit-test/webmodel"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOnDepositWallet_Unexpected_Error_When_SendMessage(t *testing.T) {
	publisherMock := pubsubMock.Publisher{}
	viewTableMock := pubsubMock.ViewTable{}

	usecase := wallet.NewWalletUsecase(wallet.UsecaseProperty{
		ServiceName:           "test-service",
		Logger:                logrus.New(),
		DepositTopicPublisher: &publisherMock,
		RollingPeriod:         180,
		Threshold:             10000,
		BalanceViewTable:      &viewTableMock,
		ThresholdViewTable:    &viewTableMock,
	})

	publisherMock.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(exception.ErrInternalServer)
	payload := webmodel.DepositWalletPayload{
		WalletId: "1",
		Amount:   1000,
	}
	resp := usecase.Deposit(context.TODO(), payload)

	assert.Error(t, resp.Error())
	assert.Equal(t, exception.ErrInternalServer, resp.Error(), "should equal to internal server error")
	assert.Equal(t, response.StatUnexpectedError, resp.Status(), "should equal to status unexpected error")
	assert.Equal(t, http.StatusInternalServerError, resp.HTTPStatusCode(), "should equal to http status internal server error")

	publisherMock.AssertExpectations(t)
	viewTableMock.AssertExpectations(t)

}

func TestOnDepositWallet_Success(t *testing.T) {
	publisherMock := pubsubMock.Publisher{}
	viewTableMock := pubsubMock.ViewTable{}

	usecase := wallet.NewWalletUsecase(wallet.UsecaseProperty{
		ServiceName:           "test-service",
		Logger:                logrus.New(),
		DepositTopicPublisher: &publisherMock,
		RollingPeriod:         180,
		Threshold:             10000,
		BalanceViewTable:      &viewTableMock,
		ThresholdViewTable:    &viewTableMock,
	})

	publisherMock.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	payload := webmodel.DepositWalletPayload{
		WalletId: "1",
		Amount:   1000,
	}
	resp := usecase.Deposit(context.TODO(), payload)

	assert.Nil(t, resp.Error())
	assert.Equal(t, response.StatOK, resp.Status(), "should equal to status ok")
	assert.Equal(t, http.StatusOK, resp.HTTPStatusCode(), "should equal to http status ok/200")

	publisherMock.AssertExpectations(t)
	viewTableMock.AssertExpectations(t)

}

func TestGetDetailWallet_UnexpectedError_When_GetFromBalanceGroupTable(t *testing.T) {
	publisherMock := pubsubMock.Publisher{}
	balanceTableMock := pubsubMock.ViewTable{}
	thresholdTableMock := pubsubMock.ViewTable{}

	usecase := wallet.NewWalletUsecase(wallet.UsecaseProperty{
		ServiceName:           "test-service",
		Logger:                logrus.New(),
		DepositTopicPublisher: &publisherMock,
		RollingPeriod:         180,
		Threshold:             10000,
		BalanceViewTable:      &balanceTableMock,
		ThresholdViewTable:    &thresholdTableMock,
	})

	balanceTableMock.On("Get", mock.AnythingOfType("string")).Return(nil, exception.ErrInternalServer)

	resp := usecase.GetDetail(context.TODO(), "1")

	assert.Error(t, resp.Error())
	assert.Equal(t, exception.ErrInternalServer, resp.Error(), "should equal to internal server error")
	assert.Equal(t, response.StatUnexpectedError, resp.Status(), "should equal to status unexpected error")
	assert.Equal(t, http.StatusInternalServerError, resp.HTTPStatusCode(), "should equal to http status internal server error")

	balanceTableMock.AssertExpectations(t)
	thresholdTableMock.AssertExpectations(t)
}

func TestGetDetailWallet_UnexpectedError_When_GetFromThresholdGroupTable(t *testing.T) {
	publisherMock := pubsubMock.Publisher{}
	balanceTableMock := pubsubMock.ViewTable{}
	thresholdTableMock := pubsubMock.ViewTable{}

	usecase := wallet.NewWalletUsecase(wallet.UsecaseProperty{
		ServiceName:           "test-service",
		Logger:                logrus.New(),
		DepositTopicPublisher: &publisherMock,
		RollingPeriod:         180,
		Threshold:             10000,
		BalanceViewTable:      &balanceTableMock,
		ThresholdViewTable:    &thresholdTableMock,
	})

	balanceTableMock.On("Get", mock.AnythingOfType("string")).Return(entity.Wallet{}, nil)
	thresholdTableMock.On("Get", mock.AnythingOfType("string")).Return(nil, exception.ErrInternalServer)

	resp := usecase.GetDetail(context.TODO(), "1")

	assert.Error(t, resp.Error())
	assert.Equal(t, exception.ErrInternalServer, resp.Error(), "should equal to internal server error")
	assert.Equal(t, response.StatUnexpectedError, resp.Status(), "should equal to status unexpected error")
	assert.Equal(t, http.StatusInternalServerError, resp.HTTPStatusCode(), "should equal to http status internal server error")

	balanceTableMock.AssertExpectations(t)
	thresholdTableMock.AssertExpectations(t)
}

func TestGetDetailWallet_NotFoundError_When_GroupTable_Empty(t *testing.T) {
	publisherMock := pubsubMock.Publisher{}
	balanceTableMock := pubsubMock.ViewTable{}
	thresholdTableMock := pubsubMock.ViewTable{}

	usecase := wallet.NewWalletUsecase(wallet.UsecaseProperty{
		ServiceName:           "test-service",
		Logger:                logrus.New(),
		DepositTopicPublisher: &publisherMock,
		RollingPeriod:         180,
		Threshold:             10000,
		BalanceViewTable:      &balanceTableMock,
		ThresholdViewTable:    &thresholdTableMock,
	})

	balanceTableMock.On("Get", mock.AnythingOfType("string")).Return(entity.Wallet{}, nil)
	thresholdTableMock.On("Get", mock.AnythingOfType("string")).Return(nil, nil)

	resp := usecase.GetDetail(context.TODO(), "1")

	assert.Error(t, resp.Error())
	assert.Equal(t, exception.ErrNotFound, resp.Error(), "should equal to not found error")
	assert.Equal(t, response.StatNotFound, resp.Status(), "should equal to status not found error")
	assert.Equal(t, http.StatusNotFound, resp.HTTPStatusCode(), "should equal to http status not found error/404")

	balanceTableMock.AssertExpectations(t)
	thresholdTableMock.AssertExpectations(t)
}

func TestGetDetailWallet_Success(t *testing.T) {
	publisherMock := pubsubMock.Publisher{}
	balanceTableMock := pubsubMock.ViewTable{}
	thresholdTableMock := pubsubMock.ViewTable{}

	usecase := wallet.NewWalletUsecase(wallet.UsecaseProperty{
		ServiceName:           "test-service",
		Logger:                logrus.New(),
		DepositTopicPublisher: &publisherMock,
		RollingPeriod:         180,
		Threshold:             10000,
		BalanceViewTable:      &balanceTableMock,
		ThresholdViewTable:    &thresholdTableMock,
	})

	wallet := &entity.Wallet{
		WalletId: "1",
		Balance:  1000,
	}
	threshold := &entity.Threshold{
		WalletId:                 "1",
		Deposit:                  1000,
		TotalDepositWithinWindow: 1000,
		StartWindowTime:          time.Now().UnixNano(),
		CreatedTime:              time.Now().UnixNano(),
		AboveThreshold:           false,
	}
	balanceTableMock.On("Get", mock.AnythingOfType("string")).Return(wallet, nil)
	thresholdTableMock.On("Get", mock.AnythingOfType("string")).Return(threshold, nil)

	resp := usecase.GetDetail(context.TODO(), "1")

	assert.Nil(t, resp.Error())
	assert.Equal(t, response.StatOK, resp.Status(), "should equal to status ok")
	assert.Equal(t, http.StatusOK, resp.HTTPStatusCode(), "should equal to http status ok/200")
	data := resp.Data().(webmodel.DetailWalletResponse)
	assert.Equal(t, data.WalletId, "1")
	assert.Equal(t, data.Balance, float64(1000))
	assert.Equal(t, data.AboveThreshold, false)

	balanceTableMock.AssertExpectations(t)
	thresholdTableMock.AssertExpectations(t)
}

func TestAddBalance_Success(t *testing.T) {
	publisherMock := pubsubMock.Publisher{}
	balanceTableMock := pubsubMock.ViewTable{}
	thresholdTableMock := pubsubMock.ViewTable{}
	contextMock := pubsubMock.GokaContext{}
	usecase := wallet.NewWalletUsecase(wallet.UsecaseProperty{
		ServiceName:           "test-service",
		Logger:                logrus.New(),
		DepositTopicPublisher: &publisherMock,
		RollingPeriod:         180,
		Threshold:             10000,
		BalanceViewTable:      &balanceTableMock,
		ThresholdViewTable:    &thresholdTableMock,
	})
	payload := &model.DepositWallet{
		WalletId: "1",
		Amount:   1000,
	}
	contextMock.On("Value").Return(nil)
	contextMock.On("SetValue", mock.Anything).Return(nil)
	resp := usecase.AddBalance(&contextMock, payload)
	assert.Nil(t, resp.Error())
	assert.Equal(t, response.StatOK, resp.Status(), "should equal to status ok")
	assert.Equal(t, http.StatusOK, resp.HTTPStatusCode(), "should equal to http status ok/200")

	contextMock.AssertExpectations(t)
}

func TestAddBalance_Success_DataAlreadyExist(t *testing.T) {
	publisherMock := pubsubMock.Publisher{}
	balanceTableMock := pubsubMock.ViewTable{}
	thresholdTableMock := pubsubMock.ViewTable{}
	contextMock := pubsubMock.GokaContext{}
	usecase := wallet.NewWalletUsecase(wallet.UsecaseProperty{
		ServiceName:           "test-service",
		Logger:                logrus.New(),
		DepositTopicPublisher: &publisherMock,
		RollingPeriod:         180,
		Threshold:             10000,
		BalanceViewTable:      &balanceTableMock,
		ThresholdViewTable:    &thresholdTableMock,
	})
	payload := &model.DepositWallet{
		WalletId: "1",
		Amount:   1000,
	}
	contextMock.On("Value").Return(&entity.Wallet{WalletId: "1", Balance: 1000})
	contextMock.On("SetValue", mock.Anything).Return(nil)
	resp := usecase.AddBalance(&contextMock, payload)
	assert.Nil(t, resp.Error())
	assert.Equal(t, response.StatOK, resp.Status(), "should equal to status ok")
	assert.Equal(t, http.StatusOK, resp.HTTPStatusCode(), "should equal to http status ok/200")
	data := resp.Data().(*entity.Wallet)
	assert.Equal(t, data.Balance, float64(2000))
	contextMock.AssertExpectations(t)
}

func TestProcessThreshold_Success(t *testing.T) {
	publisherMock := pubsubMock.Publisher{}
	balanceTableMock := pubsubMock.ViewTable{}
	thresholdTableMock := pubsubMock.ViewTable{}
	contextMock := pubsubMock.GokaContext{}
	usecase := wallet.NewWalletUsecase(wallet.UsecaseProperty{
		ServiceName:           "test-service",
		Logger:                logrus.New(),
		DepositTopicPublisher: &publisherMock,
		RollingPeriod:         180,
		Threshold:             10000,
		BalanceViewTable:      &balanceTableMock,
		ThresholdViewTable:    &thresholdTableMock,
	})
	payload := &model.DepositWallet{
		WalletId: "1",
		Amount:   1000,
	}
	contextMock.On("Value").Return(nil)
	contextMock.On("SetValue", mock.Anything).Return(nil)
	resp := usecase.ProcessThreshold(&contextMock, payload)
	assert.Nil(t, resp.Error())
	assert.Equal(t, response.StatOK, resp.Status(), "should equal to status ok")
	assert.Equal(t, http.StatusOK, resp.HTTPStatusCode(), "should equal to http status ok/200")

	contextMock.AssertExpectations(t)
}

func TestProcessThreshold_Success_DataAlreadyExist(t *testing.T) {
	publisherMock := pubsubMock.Publisher{}
	balanceTableMock := pubsubMock.ViewTable{}
	thresholdTableMock := pubsubMock.ViewTable{}
	contextMock := pubsubMock.GokaContext{}
	usecase := wallet.NewWalletUsecase(wallet.UsecaseProperty{
		ServiceName:           "test-service",
		Logger:                logrus.New(),
		DepositTopicPublisher: &publisherMock,
		RollingPeriod:         180,
		Threshold:             10000,
		BalanceViewTable:      &balanceTableMock,
		ThresholdViewTable:    &thresholdTableMock,
	})
	payload := &model.DepositWallet{
		WalletId: "1",
		Amount:   1000,
	}
	now := time.Now().UnixNano()
	threshold := &entity.Threshold{
		WalletId:                 "1",
		Deposit:                  1000,
		TotalDepositWithinWindow: 1000,
		StartWindowTime:          now,
		CreatedTime:              now,
		AboveThreshold:           false,
	}
	contextMock.On("Value").Return(threshold)
	contextMock.On("SetValue", mock.Anything).Return(nil)
	resp := usecase.ProcessThreshold(&contextMock, payload)
	assert.Nil(t, resp.Error())
	assert.Equal(t, response.StatOK, resp.Status(), "should equal to status ok")
	assert.Equal(t, http.StatusOK, resp.HTTPStatusCode(), "should equal to http status ok/200")
	data := resp.Data().(*entity.Threshold)
	assert.Equal(t, data.TotalDepositWithinWindow, float64(2000))
	contextMock.AssertExpectations(t)
}
