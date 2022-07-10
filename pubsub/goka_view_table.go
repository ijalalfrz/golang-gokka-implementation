package pubsub

import (
	"context"

	"github.com/lovoo/goka"
	"github.com/sirupsen/logrus"
)

// GokaViewTableAdapter is a concrete struct of goka kafka adapter.
type GokaViewTableAdapter struct {
	logger    *logrus.Logger
	view      *goka.View
	closeChan chan struct{}
}

func NewGokaViewTableAdapter(logger *logrus.Logger, group string, brokers []string, codec GokaCodec) ViewTable {
	v, err := goka.NewView(brokers, goka.GroupTable(goka.Group(group)), codec)
	if err != nil {
		logger.Fatal(err)
	}
	closeChan := make(chan struct{}, 1)

	view := &GokaViewTableAdapter{
		logger:    logger,
		view:      v,
		closeChan: closeChan,
	}

	return view
}

func (gk *GokaViewTableAdapter) Open() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-gk.closeChan:
			cancel()
		default:
			if err := gk.view.Run(ctx); err != nil {
				gk.logger.Errorf("Error running view: %v", err)
			} else {
				gk.logger.Info("View shutdown grafully")
			}
		}
	}()

	return
}

func (gk *GokaViewTableAdapter) Get(key string) (data interface{}, err error) {
	data, err = gk.view.Get(key)
	return
}

func (gk *GokaViewTableAdapter) Close() {
	defer close(gk.closeChan)

	gk.closeChan <- struct{}{}

	gk.logger.Info("[Goka] View Table is gracefully shut down.")
	return
}
