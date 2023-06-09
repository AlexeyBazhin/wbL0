package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/AlexeyBazhin/wbL0/internal/domain/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type (
	Server struct {
		base     *http.Server
		svc      ServerService
		logger   *zap.SugaredLogger
		ctx      context.Context
		bindAddr string
	}
	ServerService interface {
		GetOrderById(ctx context.Context, orderUid uuid.UUID) (*model.CompleteOrder, error)

		PushToCache(ctx context.Context, orderUid uuid.UUID, data []byte) error
		PullFromCache(ctx context.Context, orderUid uuid.UUID) ([]byte, error)
	}
	OptionFunc func(s *Server)
)

func New(ctx context.Context, opts ...OptionFunc) *Server {
	server := &Server{ctx: ctx}
	for _, option := range opts {
		option(server)
	}
	rtr := mux.NewRouter()
	rtr.HandleFunc("/search", server.searchHandler)
	rtr.HandleFunc("/orders/{order-uid}", server.orderHandler)

	//rtr.HandleFunc("/order", server.createOrderHandler).Methods("POST")

	server.base = &http.Server{
		Addr:    server.bindAddr,
		Handler: rtr,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
		ReadTimeout:  time.Duration(10) * time.Second,
		WriteTimeout: time.Duration(10) * time.Second,
	}

	return server
}

func (server *Server) Run() func() error {
	server.logger.Infof("[http-server] started and listening on %v", server.bindAddr)
	return func() error {
		defer server.logger.Error("[http-server] stopped")
		return server.base.ListenAndServe()
	}
}

func WithBindAddress(bindAddr string) OptionFunc {
	return func(s *Server) {
		s.bindAddr = bindAddr
	}
}
func WithService(svc ServerService) OptionFunc {
	return func(s *Server) {
		s.svc = svc
	}
}
func WithLogger(logger *zap.SugaredLogger) OptionFunc {
	return func(s *Server) {
		s.logger = logger
	}
}
