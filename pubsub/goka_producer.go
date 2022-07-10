package pubsub

import (
	"context"

	"github.com/lovoo/goka"
	"github.com/sirupsen/logrus"
)

// GokaProducerAdapter is a concrete struct of goka kafka adapter.
type GokaProducerAdapter struct {
	logger  *logrus.Logger
	emitter *goka.Emitter
}

// NewGokaProducerAdapter will create producer for produce message to kafka
func NewGokaProducerAdapter(logger *logrus.Logger, brokers []string, topic string, codec GokaCodec) (publisher Publisher, err error) {
	emitter, err := goka.NewEmitter(brokers, goka.Stream(topic), codec)
	if err != nil {
		return
	}
	publisher = &GokaProducerAdapter{
		logger:  logger,
		emitter: emitter,
	}
	return
}

// Close will close the producer
func (gk *GokaProducerAdapter) Close() (err error) {
	err = gk.emitter.Finish()
	if err == nil {
		gk.logger.Info("[Goka] Producer is gracefully shutdown")
	}
	return
}

// Send will send kafka message
func (gk *GokaProducerAdapter) Send(ctx context.Context, key string, message interface{}) (err error) {
	err = gk.emitter.EmitSync(key, message)
	return
}
