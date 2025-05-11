package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// writeError sends a structured JSON error response.
func writeError(w http.ResponseWriter, status int, errCode, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   errCode,
		Message: message,
	})
}

// NewLedger creates and returns a new in-memory ledger.
func NewLedger() *Ledger {
	return &Ledger{
		ID:                 uuid.New(),
		TransactionHistory: make([]Transaction, 0),
		Balance:            0.00,
		NextID:             1,
	}
}

// CreateTransaction handles POST /transactions.
func (l *Ledger) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var body struct {
		Type   string  `json:"type"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Invalid JSON payload")
		return
	}

	if body.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "invalid_amount", "Amount must be greater than zero")
		return
	}

	var txType TransactionType
	switch body.Type {
	case string(Deposit):
		txType = Deposit
		l.Balance += body.Amount
	case string(Withdrawal):
		if l.Balance < body.Amount {
			writeError(w, http.StatusBadRequest, "insufficient_funds", "Not enough balance to complete withdrawal")
			return
		}
		txType = Withdrawal
		l.Balance -= body.Amount
	default:
		writeError(w, http.StatusBadRequest, "invalid_transaction_type", "Transaction type must be 'deposit' or 'withdrawal'")
		return
	}

	tx := Transaction{
		ID:        l.NextID,
		Type:      txType,
		Amount:    body.Amount,
		Timestamp: time.Now(),
	}
	l.TransactionHistory = append(l.TransactionHistory, tx)
	l.NextID++

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}

// GetBalance handles GET /balance.
func (l *Ledger) GetBalance(w http.ResponseWriter, r *http.Request) {
	l.mu.Lock()
	defer l.mu.Unlock()

	resp := BalanceResponse{
		Balance:   l.Balance,
		CheckedAt: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetTransactionHistory handles GET /transactions.
func (l *Ledger) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	l.mu.Lock()
	defer l.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(l.TransactionHistory)
}
