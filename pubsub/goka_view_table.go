package pubsub

import (
	"context"

	"github.com/lovoo/goka"
	"github.com/sirupsen/logrus"
)

// GokaViewTableAdapter is a concrete struct of goka group table adapter.
type GokaViewTableAdapter struct {
	logger *logrus.Logger
	view   *goka.View
	cancel context.CancelFunc
}

// NewGokaViewTableAdapter will create goka view
func NewGokaViewTableAdapter(logger *logrus.Logger, group string, brokers []string, codec GokaCodec) (view ViewTable, err error) {
	v, err := goka.NewView(brokers, goka.GroupTable(goka.Group(group)), codec)
	if err != nil {
		return
	}

	view = &GokaViewTableAdapter{
		logger: logger,
		view:   v,
	}

	return
}

// Open will run the view on new goroutine
func (gk *GokaViewTableAdapter) Open() {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		if gk.view.CurrentState() != goka.ViewStateRunning {
			if err := gk.view.Run(ctx); err != nil {
				gk.logger.Errorf("Error running view: %v", err)
			}
		}
	}(ctx)
	gk.cancel = cancel
	return

}

// Get will return data from group table based on key
func (gk *GokaViewTableAdapter) Get(key string) (data interface{}, err error) {
	data, err = gk.view.Get(key)
	return
}

// Close will cancel view context
func (gk *GokaViewTableAdapter) Close() {
	gk.cancel()
	gk.logger.Info("[Goka] View Table is gracefully shut down.")
	return
}
