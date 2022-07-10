package pubsub

import (
	"context"

	"github.com/lovoo/goka"
	"github.com/sirupsen/logrus"
)

type GokaConsumserGroupAdapter struct {
	logger    *logrus.Logger
	processor *goka.Processor
}

// NewGokaConsumerGroupFullConfigAdapter will create consumer group and group table
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
		return
	}
	subscriber = &GokaConsumserGroupAdapter{
		logger:    logger,
		processor: p,
	}

	return
}

// Subscribe will consume the published message
func (gk *GokaConsumserGroupAdapter) Subscribe() {
	ctx := context.Background()
	go func() {
		if err := gk.processor.Run(ctx); err != nil {
			gk.logger.Errorf("Error running processor: %v", err)
		}
	}()

	return
}

// Close will stop the kafka consumer
func (gk *GokaConsumserGroupAdapter) Close() (err error) {
	gk.processor.Stop()
	gk.logger.Info("[Goka] Consumer is gracefully shut down.")
	return
}
