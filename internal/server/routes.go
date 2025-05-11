package server

import (
	"net/http"

	"github.com/elemhunt/tiny_ledger/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// loadRoutes initializes the main router and mounts all application routes.
// Returns a pointer to a chi.Mux router.
func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)

	// Health check or root route
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Ledger-specific routes
	router.Route("/ledger", loadLedgerRoutes)

	return router
}

// loadLedgerRoutes mounts all ledger-related HTTP endpoints onto the given router.
// Params:
//   - router: A chi.Router instance onto which the ledger endpoints are registered.
func loadLedgerRoutes(router chi.Router) {
	ledgerHandler := handler.NewLedger()

	router.Post("/transactions", ledgerHandler.CreateTransaction)
	router.Get("/balance", ledgerHandler.GetBalance)
	router.Get("/transaction_history", ledgerHandler.GetTransactionHistory)
}
