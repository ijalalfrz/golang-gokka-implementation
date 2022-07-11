package wallet

import (
	"fmt"

	"github.com/ijalalfrz/coinbit-test/model"
	"github.com/ijalalfrz/coinbit-test/pubsub"
	"google.golang.org/protobuf/proto"
)

type depositCodec struct {
}

func NewDepositCodec() pubsub.GokaCodec {
	return &depositCodec{}
}

func (jc *depositCodec) Encode(value interface{}) ([]byte, error) {
	if _, isDeposit := value.(*model.DepositWallet); !isDeposit {
		return nil, fmt.Errorf("Codec requires value *model.DepositWallet, got %T", value)
	}
	v := value.(*model.DepositWallet)
	return proto.Marshal(v)
}

// Decodes a deposit from []byte to it's go representation.
func (jc *depositCodec) Decode(data []byte) (interface{}, error) {
	var (
		deposit model.DepositWallet
		err     error
	)
	err = proto.Unmarshal(data, &deposit)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling deposit: %v", err)
	}
	return &deposit, nil
}
