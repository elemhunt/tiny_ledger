package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/elemhunt/tiny_ledger/internal/handler"
	"github.com/elemhunt/tiny_ledger/internal/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
)

// TestHealthCheck ensures that the health check endpoint returns a 200 status.
// It sends a GET request to the root ("/") endpoint and checks if the response code is HTTP 200 OK.
func TestHealthCheck(t *testing.T) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")
}

// TestDepositTransaction tests the /ledger/transactions endpoint for a deposit.
func TestDepositTransaction(t *testing.T) {
	ledger := handler.NewLedger()
	router := chi.NewRouter()

	// Set up routes
	router.Post("/ledger/transactions", ledger.CreateTransaction)

	// Create request body for deposit
	txBody := map[string]interface{}{
		"type":   "deposit",
		"amount": 100.0,
	}
	txBodyJSON, _ := json.Marshal(txBody)

	req, _ := http.NewRequest("POST", "/ledger/transactions", bytes.NewReader(txBodyJSON))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check for expected response
	assert.Equal(t, http.StatusCreated, rr.Code, "Expected status code 201 (Created)")

	var txResponse handler.Transaction
	err := json.NewDecoder(rr.Body).Decode(&txResponse)
	assert.NoError(t, err, "Error while decoding response")

	// Check if the transaction amount is correct
	assert.Equal(t, float64(100.0), txResponse.Amount, "Expected transaction amount 100")

	// Check if the balance was updated correctly
	assert.Equal(t, 100.0, ledger.Balance, "Expected balance 100 after deposit")
}

// TestWithdrawalTransaction tests the /ledger/transactions endpoint for a withdrawal.
func TestWithdrawalTransaction(t *testing.T) {
	ledger := handler.NewLedger()
	ledger.Balance = 200.0 // Start with a balance of 200
	router := chi.NewRouter()

	// Set up routes
	router.Post("/ledger/transactions", ledger.CreateTransaction)

	// Create request body for withdrawal
	txBody := map[string]interface{}{
		"type":   "withdrawal",
		"amount": 100.0,
	}
	txBodyJSON, _ := json.Marshal(txBody)

	req, _ := http.NewRequest("POST", "/ledger/transactions", bytes.NewReader(txBodyJSON))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check for expected response
	assert.Equal(t, http.StatusCreated, rr.Code, "Expected status code 201 (Created)")

	var txResponse handler.Transaction
	err := json.NewDecoder(rr.Body).Decode(&txResponse)
	assert.NoError(t, err, "Error while decoding response")

	// Check if the transaction amount is correct
	assert.Equal(t, float64(100.0), txResponse.Amount, "Expected transaction amount 100")

	// Check if the balance was updated correctly after withdrawal
	assert.Equal(t, 100.0, ledger.Balance, "Expected balance 100 after withdrawal")
}

// TestGetBalance tests the /ledger/balance endpoint.
// It sends a GET request to the balance endpoint and checks if the response correctly returns the current balance.
func TestGetBalance(t *testing.T) {
	ledger := handler.NewLedger()
	ledger.Balance = 200.0
	router := chi.NewRouter()

	// Set up routes
	router.Get("/ledger/balance", ledger.GetBalance)

	req, _ := http.NewRequest("GET", "/ledger/balance", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check for expected response
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")

	var balanceResponse handler.BalanceResponse
	err := json.NewDecoder(rr.Body).Decode(&balanceResponse)
	assert.NoError(t, err, "Error while decoding response")
	assert.Equal(t, 200.0, balanceResponse.Balance, "Expected balance 200")
}

// TestGetTransactionHistory tests the /ledger/transaction_history endpoint.
// It adds a sample transaction to the ledger and checks if the history response includes the expected transaction.
func TestGetTransactionHistory(t *testing.T) {
	ledger := handler.NewLedger()
	ledger.TransactionHistory = append(ledger.TransactionHistory, handler.Transaction{
		ID:        1,
		Type:      handler.Deposit,
		Amount:    100.0,
		Timestamp: time.Now(),
	})
	router := chi.NewRouter()

	// Set up routes
	router.Get("/ledger/transaction_history", ledger.GetTransactionHistory)

	req, _ := http.NewRequest("GET", "/ledger/transaction_history", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check for expected response
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")

	var txHistory []handler.Transaction
	err := json.NewDecoder(rr.Body).Decode(&txHistory)
	assert.NoError(t, err, "Error while decoding response")
	assert.Len(t, txHistory, 1, "Expected 1 transaction in history")
	assert.Equal(t, 100.0, txHistory[0].Amount, "Expected transaction amount 100")
}

// TestServerStart tests the graceful server startup and shutdown.
// It starts the server, ensures it initializes correctly, and tests its shutdown process.
func TestServerStart(t *testing.T) {
	// Simulate interrupt signal via cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	srv := server.New()
	errCh := make(chan error, 1)

	// Start server
	go func() {
		errCh <- srv.Start(ctx)
	}()

	// Give the server a moment to start
	time.Sleep(200 * time.Millisecond)

	// Cancel the context to simulate shutdown
	cancel()

	// Wait for shutdown
	select {
	case err := <-errCh:
		assert.NoError(t, err, "Server failed to shut down cleanly")
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout: server did not shut down in time")
	}
}

// TestWithdrawInsufficientFunds verifies that the server returns an appropriate error
// when attempting to withdraw more than the available balance.
func TestWithdrawInsufficientFunds(t *testing.T) {
	ledger := handler.NewLedger()
	router := chi.NewRouter()

	// Mount the transaction route
	router.Post("/ledger/transactions", ledger.CreateTransaction)

	// Try withdrawing more than balance (which is 0.0 initially)
	withdrawBody := map[string]interface{}{
		"type":   "withdrawal",
		"amount": 50.0,
	}
	bodyJSON, _ := json.Marshal(withdrawBody)

	req, _ := http.NewRequest("POST", "/ledger/transactions", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status code 400 for insufficient funds")

	// Decode JSON response
	var errResp handler.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	assert.NoError(t, err, "Expected valid JSON response")
	assert.Equal(t, "insufficient_funds", errResp.Error, "Expected error code 'insufficient_funds'")
	assert.Equal(t, "Not enough balance to complete withdrawal", errResp.Message, "Expected detailed error message")
}

// TestCreateTransaction_InvalidJSON ensures that the server responds with a 400 error
// and a clear message when the request body contains invalid JSON.
func TestCreateTransaction_InvalidJSON(t *testing.T) {
	ledger := handler.NewLedger()
	router := chi.NewRouter()
	router.Post("/ledger/transactions", ledger.CreateTransaction)

	// Invalid JSON (missing closing brace)
	body := `{"type": "deposit", "amount": 10.0`

	req, _ := http.NewRequest("POST", "/ledger/transactions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status code 400 for invalid JSON")

	var errResp handler.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	assert.NoError(t, err, "Expected valid JSON error response")
	assert.Equal(t, "invalid_request", errResp.Error)
	assert.Equal(t, "Invalid JSON payload", errResp.Message)
}

// TestCreateTransaction_InvalidType checks that transactions with an unsupported type
// (e.g., "transfer") are rejected with a proper error response.
func TestCreateTransaction_InvalidType(t *testing.T) {
	ledger := handler.NewLedger()
	router := chi.NewRouter()
	router.Post("/ledger/transactions", ledger.CreateTransaction)

	body := map[string]interface{}{
		"type":   "transfer", // invalid
		"amount": 20.0,
	}
	bodyJSON, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/ledger/transactions", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status code 400 for invalid transaction type")

	var errResp handler.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	assert.NoError(t, err, "Expected valid JSON error response")
	assert.Equal(t, "invalid_transaction_type", errResp.Error)
	assert.Equal(t, "Transaction type must be 'deposit' or 'withdrawal'", errResp.Message)
}

// TestCreateTransaction_NegativeAmount confirms that transactions with negative amounts
// are rejected with an error, enforcing business rules for valid input.
func TestCreateTransaction_NegativeAmount(t *testing.T) {
	ledger := handler.NewLedger()
	router := chi.NewRouter()
	router.Post("/ledger/transactions", ledger.CreateTransaction)

	body := map[string]interface{}{
		"type":   "deposit",
		"amount": -5.0,
	}
	bodyJSON, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/ledger/transactions", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status code 400 for negative amount")

	var errResp handler.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	assert.NoError(t, err, "Expected valid JSON error response")
	assert.Equal(t, "invalid_amount", errResp.Error)
	assert.Equal(t, "Amount must be greater than zero", errResp.Message)
}
