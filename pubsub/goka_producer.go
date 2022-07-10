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

func NewGokaProducerAdapter(logger *logrus.Logger, brokers []string, topic string, codec GokaCodec) Publisher {
	emitter, err := goka.NewEmitter(brokers, goka.Stream(topic), codec)
	if err != nil {
		logger.Fatal(err)
	}
	publisher := &GokaProducerAdapter{
		logger:  logger,
		emitter: emitter,
	}
	return publisher
}

func (gk *GokaProducerAdapter) Close() (err error) {
	err = gk.emitter.Finish()
	if err == nil {
		gk.logger.Info("[Goka] Producer is gracefully shutdown")
	}
	return
}

func (gk *GokaProducerAdapter) Send(ctx context.Context, key string, message interface{}) (err error) {
	err = gk.emitter.EmitSync(key, message)
	return
}
