package apiutils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func Serve(
	routes http.Handler,
	logger *slog.Logger,
	env string,
	port int,
	version string,
	ctx context.Context,
) error {
	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", port),
		Handler:  routes,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error)

	go func() {
		<-ctx.Done()
		logger.Info("shutting down server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			errCh <- fmt.Errorf("failed to shutdown http server: %w", err)
		}
		errCh <- nil
	}()

	logger.Info("starting server", "addr", srv.Addr, "env", env, "version", version)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		cancel()
		<-errCh
		return err
	}

	err = <-errCh
	if err != nil {
		return err
	}

	logger.Info("stopped server", "addr", srv.Addr)

	return nil
}
