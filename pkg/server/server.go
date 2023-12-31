package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type Server struct {
	http     *http.Server
	grpc     *grpc.Server
	listener net.Listener
}

type Configuration func(r *Server) error

func New(configs ...Configuration) (r *Server, err error) {
	// создаем сервер
	r = &Server{}

	for _, cfg := range configs {
		if err = cfg(r); err != nil {
			return
		}
	}
	return
}

func (s *Server) Run(logger *zap.Logger) (err error) {
	if s.http != nil {
		go func() {
			if err = s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error("ERR_INIT_REST", zap.Error(err))
				return
			}
		}()
		logger.Info("http server started on http://localhost" + s.http.Addr)
	}

	if s.grpc != nil {
		if err = s.grpc.Serve(s.listener); err != nil {
			return
		}
		logger.Info("grpc server started on http://localhost" + s.listener.Addr().String())
	}

	return
}

func (s *Server) Stop(ctx context.Context) (err error) {
	if s.http != nil {
		if err = s.http.Shutdown(ctx); err != nil {
			return
		}
	}

	return
}

func WithGRPCServer(port string) Configuration {
	return func(s *Server) (err error) {
		s.listener, err = net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
		if err != nil {
			return
		}
		s.grpc = &grpc.Server{}

		return
	}
}

func WithHTTPServer(handler http.Handler, port string) Configuration {
	return func(s *Server) (err error) {
		s.http = &http.Server{
			Handler: handler,
			Addr:    ":" + port,
		}
		return
	}
}
