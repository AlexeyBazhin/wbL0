package stanListener

import (
	"context"

	"github.com/AlexeyBazhin/wbL0/internal/api"
	"github.com/AlexeyBazhin/wbL0/internal/domain/model"
	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type (
	StanListener struct {
		svc      StanService
		stanConn stan.Conn
		subject  string
		ctx      context.Context
		logger   *zap.SugaredLogger
	}
	StanService interface {
		CreateOrder(ctx context.Context, order api.OrderJSON) (*model.Order, error)
		CreateDelivery(ctx context.Context, delivery api.DeliveryJSON, orderUid uuid.UUID) (*model.Delivery, error)
		CreatePayment(ctx context.Context, payment api.PaymentJSON, orderUid uuid.UUID) (*model.Payment, error)
		CreateItem(ctx context.Context, item api.ItemJSON, orderUid uuid.UUID) (*model.Item, error)
		InsertCompleteOrder(ctx context.Context, completeModels *model.CompleteOrder) error

		PushToCache(ctx context.Context, orderUid uuid.UUID, data []byte) error
	}
	OptionFunc func(apiStan *StanListener)
)

func New(opts ...OptionFunc) *StanListener {
	apiStan := new(StanListener)
	for _, option := range opts {
		option(apiStan)
	}

	return apiStan
}

func (stanListener *StanListener) Run() func() error {
	return func() error {
		_, err := stanListener.stanConn.Subscribe(stanListener.subject, stanListener.stanHandler, stan.DurableName("wbL0"))
		if err != nil {
			stanListener.logger.Errorf("[nats-streaming] cannot subscribe: %w", err)
			return err
		}
		stanListener.logger.Info("[nats-streaming] successfuly subscribed")
		//TODO
		//sub.Close()
		return nil
	}

}
func WithLogger(logger *zap.SugaredLogger) OptionFunc {
	return func(stanListener *StanListener) {
		stanListener.logger = logger
	}
}
func WithContext(ctx context.Context) OptionFunc {
	return func(stanListener *StanListener) {
		stanListener.ctx = ctx
	}
}
func WithService(svc StanService) OptionFunc {
	return func(stanListener *StanListener) {
		stanListener.svc = svc
	}
}
func WithStanConn(stanConn stan.Conn) OptionFunc {
	return func(stanListener *StanListener) {
		stanListener.stanConn = stanConn
	}
}
func WithSubject(subject string) OptionFunc {
	return func(stanListener *StanListener) {
		stanListener.subject = subject
	}
}
