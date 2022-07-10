package wallet

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ijalalfrz/coinbit-test/entity"
	"github.com/ijalalfrz/coinbit-test/model"
	"github.com/ijalalfrz/coinbit-test/pubsub"
	"github.com/ijalalfrz/coinbit-test/response"
	"github.com/ijalalfrz/coinbit-test/webmodel"
	"github.com/lovoo/goka"
	"github.com/sirupsen/logrus"
)

// collection of message
const (
	depositUnexpectedErrMessage    = "Unexpected error while processing deposit wallet"
	depositSuccessMessage          = "Deposit to wallet has been processed"
	addBalanceSuccessMessage       = "Add balance to wallet: %s is successfully processed, current balance: %.2f"
	processThresholdSuccessMessage = "Balance threshold for wallet: %s has been processed, current above threshold status: %t"
	detailUnexpectedErrMessage     = "Unexpected error while getting wallet details"
	detailSuccessMessage           = "Detail wallet"
)

// Usecase is a collection of behavior of wallet.
type Usecase interface {
	Deposit(ctx context.Context, payload webmodel.DepositWalletPayload) (resp response.Response)
	AddBalance(ctx goka.Context, payload *model.DepositWallet) (resp response.Response)
	ProcessThreshold(ctx goka.Context, payload *model.DepositWallet) (resp response.Response)
	GetDetail(ctx context.Context, walletId string) (resp response.Response)
}

type walletUsecase struct {
	serviceName           string
	logger                *logrus.Logger
	depositTopicPublisher pubsub.Publisher
	rollingPeriod         int
	threshold             int64
	balanceViewTable      pubsub.ViewTable
	thresholdViewTable    pubsub.ViewTable
}

func NewWalletUsecase(property UsecaseProperty) Usecase {
	return &walletUsecase{
		serviceName:           property.ServiceName,
		logger:                property.Logger,
		depositTopicPublisher: property.DepositTopicPublisher,
		rollingPeriod:         property.RollingPeriod,
		threshold:             property.Threshold,
		balanceViewTable:      property.BalanceViewTable,
		thresholdViewTable:    property.ThresholdViewTable,
	}
}

// Deposit is a method for add balance to wallet
func (u walletUsecase) Deposit(ctx context.Context, payload webmodel.DepositWalletPayload) (resp response.Response) {
	var deposit = &model.DepositWallet{
		WalletId: payload.WalletId,
		Amount:   payload.Amount,
	}

	err := u.depositTopicPublisher.Send(ctx, payload.WalletId, deposit)
	if err != nil {
		u.logger.Error(err)
		return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, depositUnexpectedErrMessage)
	}

	return response.NewSuccessResponse(nil, response.StatOK, depositSuccessMessage)
}

func (u walletUsecase) AddBalance(ctx goka.Context, payload *model.DepositWallet) (resp response.Response) {
	var wallet *entity.Wallet
	if val := ctx.Value(); val != nil {
		wallet = val.(*entity.Wallet)
	} else {
		wallet = new(entity.Wallet)
	}

	wallet.Balance += payload.GetAmount()
	wallet.WalletId = payload.GetWalletId()
	ctx.SetValue(wallet)
	return response.NewSuccessResponse(nil, response.StatOK, fmt.Sprintf(addBalanceSuccessMessage, wallet.WalletId, wallet.Balance))

}

func (u walletUsecase) ProcessThreshold(ctx goka.Context, payload *model.DepositWallet) (resp response.Response) {
	var threshold *entity.Threshold
	if val := ctx.Value(); val != nil {
		threshold = val.(*entity.Threshold)
	} else {
		threshold = new(entity.Threshold)
		threshold.StartWindowTime = time.Now().UnixNano()
	}

	now := time.Now().UnixNano()
	threshold.WalletId = payload.GetWalletId()
	threshold.Deposit = payload.Amount
	threshold.TotalDepositWithinWindow += payload.Amount
	threshold.CreatedTime = now

	timeNow := time.Unix(0, now)
	timeStartRollingPeriod := time.Unix(0, threshold.StartWindowTime)
	diff := timeNow.Sub(timeStartRollingPeriod)
	if diff.Seconds() > float64(u.rollingPeriod) {
		threshold.AboveThreshold = false
		threshold.StartWindowTime = now
		threshold.TotalDepositWithinWindow = payload.Amount
	} else {
		if threshold.TotalDepositWithinWindow > float64(u.threshold) {
			threshold.AboveThreshold = true
		} else {
			threshold.AboveThreshold = false
		}
	}
	ctx.SetValue(threshold)
	return response.NewSuccessResponse(nil, response.StatOK, fmt.Sprintf(processThresholdSuccessMessage, threshold.WalletId, threshold.AboveThreshold))

}

func (u walletUsecase) GetDetail(ctx context.Context, walletId string) (resp response.Response) {
	balanceData, err := u.balanceViewTable.Get(walletId)
	if err != nil {
		u.logger.Error(err)
		return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, detailUnexpectedErrMessage)
	}

	thresholdData, err := u.thresholdViewTable.Get(walletId)
	if err != nil {
		u.logger.Error(err)
		return response.NewErrorResponse(err, http.StatusInternalServerError, nil, response.StatUnexpectedError, detailUnexpectedErrMessage)
	}

	balance := balanceData.(*entity.Wallet)
	threshold := thresholdData.(*entity.Threshold)
	detail := webmodel.DetailWalletResponse{
		WalletId:       balance.WalletId,
		Balance:        balance.Balance,
		AboveThreshold: threshold.AboveThreshold,
	}
	return response.NewSuccessResponse(detail, response.StatOK, detailSuccessMessage)
}
