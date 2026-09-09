package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func Run(ctx context.Context, port int16, handler http.Handler) error {
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(int(port)),
		Handler: handler,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpServer.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
