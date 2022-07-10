package pubsub

import (
	"context"

	"github.com/lovoo/goka"
	"github.com/sirupsen/logrus"
)

type GokaCodec interface {
	Encode(value interface{}) ([]byte, error)
	Decode(data []byte) (interface{}, error)
}
type GokaConsumserGroupAdapter struct {
	logger    *logrus.Logger
	processor *goka.Processor
	closeChan chan struct{}
}

func NewGokaConsumerGroupFullConfigAdapter(
	logger *logrus.Logger, addresses []string, groupID string, topic string, handler GokaEventHandler,
	topicManagerConfig *goka.TopicManagerConfig, inputCodec GokaCodec, tableCodec GokaCodec,
) (subscriber Subscriber, err error) {
	g := goka.DefineGroup(goka.Group(groupID),
		goka.Input(goka.Stream(topic), inputCodec, handler.Handle),
		goka.Persist(tableCodec),
	)
	p, err := goka.NewProcessor(addresses,
		g,
		goka.WithTopicManagerBuilder(goka.TopicManagerBuilderWithTopicManagerConfig(topicManagerConfig)),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)
	if err != nil {
		logger.Fatal(err)
	}
	closeChan := make(chan struct{}, 1)

	subscriber = &GokaConsumserGroupAdapter{
		logger:    logger,
		processor: p,
		closeChan: closeChan,
	}

	return
}

// Subscribe will consume the published message
func (gk *GokaConsumserGroupAdapter) Subscribe() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-gk.closeChan:
			cancel()
		default:
			if err := gk.processor.Run(ctx); err != nil {
				gk.logger.Errorf("Error running processor: %v", err)
			} else {
				gk.logger.Info("Processor shutdown grafully")
			}
		}
	}()

	return
}

// Close will stop the kafka consumer
func (gk *GokaConsumserGroupAdapter) Close() (err error) {
	defer close(gk.closeChan)

	gk.closeChan <- struct{}{}

	gk.processor.Stop()

	gk.logger.Info("[Goka] Consumer is gracefully shut down.")
	return
}
