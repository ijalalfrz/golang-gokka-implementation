package wallet

import (
	"encoding/json"
	"fmt"

	"github.com/ijalalfrz/coinbit-test/entity"
	"github.com/ijalalfrz/coinbit-test/pubsub"
)

type threshold struct {
}

func NewThresholdCodec() pubsub.GokaCodec {
	return &threshold{}
}

func (c *threshold) Encode(value interface{}) ([]byte, error) {
	if _, isDeposit := value.(*entity.Threshold); !isDeposit {
		return nil, fmt.Errorf("Codec requires value *entity.Threshold, got %T", value)
	}
	v := value.(*entity.Threshold)
	return json.Marshal(v)
}

// Decodes a user from []byte to it's go representation.
func (c *threshold) Decode(data []byte) (interface{}, error) {
	var (
		threshold entity.Threshold
		err       error
	)
	err = json.Unmarshal(data, &threshold)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling threshold: %v", err)
	}
	return &threshold, nil
}
