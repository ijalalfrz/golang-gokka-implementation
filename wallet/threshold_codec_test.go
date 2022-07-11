package wallet_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ijalalfrz/coinbit-test/entity"
	"github.com/ijalalfrz/coinbit-test/model"
	"github.com/ijalalfrz/coinbit-test/wallet"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestThresholdCodec_Error_Encode(t *testing.T) {
	codec := wallet.NewThresholdCodec()

	result, err := codec.Encode(nil)

	assert.Nil(t, result, "should be null")
	assert.Error(t, err, "should be error")
}

func TestThresholdCodec_Success_Encode(t *testing.T) {
	codec := wallet.NewThresholdCodec()

	data := &entity.Threshold{
		WalletId:                 "1",
		Deposit:                  1000,
		TotalDepositWithinWindow: 1000,
		StartWindowTime:          time.Now().UnixNano(),
		CreatedTime:              time.Now().UnixNano(),
		AboveThreshold:           false,
	}
	result, err := codec.Encode(data)

	assert.Nil(t, err, "should be null")
	assert.NotNil(t, result, "should be not null")
}

func TestThresholdCodec_Success_Decode(t *testing.T) {
	codec := wallet.NewThresholdCodec()
	data := &entity.Threshold{
		WalletId: "1",
	}
	buff, err := json.Marshal(data)
	result, err := codec.Decode(buff)

	assert.Nil(t, err, "should be null")
	res := result.(*entity.Threshold)
	assert.Equal(t, res.WalletId, "1")
}

func TestThresholdCodec_Error_Decode(t *testing.T) {
	codec := wallet.NewThresholdCodec()
	data := &model.DepositWallet{
		WalletId: "1",
	}
	buff, err := proto.Marshal(data)
	result, err := codec.Decode(buff)

	assert.Error(t, err, "should be error")
	assert.Nil(t, result, "should be null")
}
