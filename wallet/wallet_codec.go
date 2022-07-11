package wallet

import (
	"encoding/json"
	"fmt"

	"github.com/ijalalfrz/coinbit-test/entity"
	"github.com/ijalalfrz/coinbit-test/pubsub"
)

type walletCodec struct {
}

func NewWalletCodec() pubsub.GokaCodec {
	return &walletCodec{}
}

func (jc *walletCodec) Encode(value interface{}) ([]byte, error) {
	if _, isDeposit := value.(*entity.Wallet); !isDeposit {
		return nil, fmt.Errorf("Codec requires value *entity.Wallet, got %T", value)
	}
	v := value.(*entity.Wallet)
	return json.Marshal(v)
}

// Decodes a wallet from []byte to it's go representation.
func (jc *walletCodec) Decode(data []byte) (interface{}, error) {
	var (
		deposit entity.Wallet
		err     error
	)
	err = json.Unmarshal(data, &deposit)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling deposit: %v", err)
	}
	return &deposit, nil
}
