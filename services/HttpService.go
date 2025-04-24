package services

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type HttpService struct {
	server   *http.Server
	initFunc func(*http.Server, *ServiceProvider) error
	sp       *ServiceProvider
}

func NewHttpService(server *http.Server, initFunc func(httpServer *http.Server, serviceProvider *ServiceProvider) error) *HttpService {
	return &HttpService{
		server:   server,
		initFunc: initFunc,
	}
}

func (h *HttpService) Init(provider *ServiceProvider) error {
	h.sp = provider
	if h.initFunc != nil {
		return h.initFunc(h.server, provider)
	}
	return nil
}

// Run starts the http server and waits for it to be stopped.
// This returns the error of net/http.Server.ListenAndServe(), so when error is not nil, it can be ErrServerClosed or others.
func (h *HttpService) Run(ctx context.Context) error {
	// register ctx for shutdown
	var stoppingCtx context.Context
	var stoppingCancel context.CancelFunc
	go func() {
		<-ctx.Done()
		stoppingCtx, stoppingCancel = context.WithTimeout(context.Background(), 30*time.Second)
		go func() {
			time.Sleep(time.Millisecond * 10) // give time to the shutdown to start
			select {
			case <-stoppingCtx.Done():
				break
			default:
			}
			time.Sleep(time.Second * 2)
			select {
			case <-stoppingCtx.Done():
				break
			default:
				h.sp.Logger.Information("Waiting for server to shutdown")
			}
			time.Sleep(time.Second * 10)
			select {
			case <-stoppingCtx.Done():
				break
			default:
				h.sp.Logger.Warning("Service is still running after 10 seconds, waiting for another 20 seconds before force shutdown")
			}
		}()
		err := h.server.Shutdown(stoppingCtx)
		stoppingCancel()
		switch {
		case err != nil && errors.Is(err, context.DeadlineExceeded):
			h.sp.Logger.Warning("Deadline exceeded, service is killed")
			break
		case err != nil && errors.Is(err, http.ErrServerClosed):
		case err == nil:
			h.sp.Logger.Information("Service stopped")
			break
		case err != nil:
			h.sp.Logger.Warning("Error stopping service: %v", err)
			break
		}
	}()
	serverErr := h.server.ListenAndServe()
	<-stoppingCtx.Done()
	return serverErr
}
