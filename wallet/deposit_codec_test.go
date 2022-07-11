package wallet_test

import (
	"encoding/json"
	"testing"

	"github.com/ijalalfrz/coinbit-test/entity"
	"github.com/ijalalfrz/coinbit-test/model"
	"github.com/ijalalfrz/coinbit-test/wallet"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestDepositCodec_Error_Encode(t *testing.T) {
	codec := wallet.NewDepositCodec()

	result, err := codec.Encode(nil)

	assert.Nil(t, result, "should be null")
	assert.Error(t, err, "should be error")
}

func TestDepositCodec_Success_Encode(t *testing.T) {
	codec := wallet.NewDepositCodec()

	data := &model.DepositWallet{
		WalletId: "1",
		Amount:   1000,
	}
	result, err := codec.Encode(data)

	assert.Nil(t, err, "should be null")
	assert.NotNil(t, result, "should be not null")
}

func TestDepositCodec_Success_Decode(t *testing.T) {
	codec := wallet.NewDepositCodec()
	data := &model.DepositWallet{
		WalletId: "1",
		Amount:   100,
	}
	buff, err := proto.Marshal(data)
	result, err := codec.Decode(buff)

	assert.Nil(t, err, "should be null")
	res := result.(*model.DepositWallet)
	assert.Equal(t, res.WalletId, "1")
}

func TestDepositCodec_Error_Decode(t *testing.T) {
	codec := wallet.NewDepositCodec()
	data := &entity.Wallet{
		WalletId: "1",
	}
	buff, err := json.Marshal(data)
	result, err := codec.Decode(buff)

	assert.Error(t, err, "should be error")
	assert.Nil(t, result, "should be null")
}
