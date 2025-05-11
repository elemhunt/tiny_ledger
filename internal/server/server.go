package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/elemhunt/tiny_ledger/config"
	"github.com/elemhunt/tiny_ledger/internal/handler"
)

type Server struct {
	// Struct to store server dependencies
	router http.Handler
	ledger *handler.Ledger
}

// New creates a new Server instance.
func New() *Server {
	ledger := handler.NewLedger()

	return &Server{
		router: loadRoutes(),
		ledger: ledger,
	}
}

// Start launches the server and handles graceful shutdown.
func (s *Server) Start(ctx context.Context) error {
	config.LoadEnv()
	port := config.GetEnv("PORT", "8080")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: s.router,
	}

	fmt.Println("Starting server on PORT:", port, "...")

	// Set error channel
	errCh := make(chan error, 1)

	// Run server concurrently
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			//	Send error to error channel and close channel
			errCh <- fmt.Errorf("failed to start server: %w", err)
		}
		close(errCh)
	}()

	//	Select case for desiding onhow to close with graceful shutdown
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		// Shutdown gracefully with a timeout
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(timeoutCtx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
		return nil
	}
}
