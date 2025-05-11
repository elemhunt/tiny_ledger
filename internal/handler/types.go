package handler

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// TransactionType represents the type of a ledger transaction.
type TransactionType string

const (
	// Deposit indicates a deposit transaction.
	Deposit TransactionType = "deposit"

	// Withdrawal indicates a withdrawal transaction.
	Withdrawal TransactionType = "withdrawal"
)

// LedgerHandler wraps the service layer (optional, for potential expansion).
type LedgerHandler struct {
	Service *Ledger
}

// Transaction represents a single ledger entry.
type Transaction struct {
	ID        int64           `json:"id"`        // Unique transaction identifier.
	Type      TransactionType `json:"type"`      // Transaction type: "deposit" or "withdrawal".
	Amount    float64         `json:"amount"`    // Transaction amount.
	Timestamp time.Time       `json:"timestamp"` // Time the transaction occurred.
}

// BalanceResponse is the response returned by the balance endpoint.
type BalanceResponse struct {
	Balance   float64   `json:"balance"`    // Current account balance.
	CheckedAt time.Time `json:"checked_at"` // Timestamp when the balance was retrieved.
}

// Ledger represents the entire state of the in-memory ledger.
type Ledger struct {
	ID                 uuid.UUID     `json:"id"`                  // Unique identifier for the ledger instance.
	TransactionHistory []Transaction `json:"transaction_history"` // List of all transactions in the ledger.
	Balance            float64       `json:"balance"`             // Current balance in the ledger.
	NextID             int64         `json:"next_id"`             // Next transaction ID to be assigned.

	mu sync.Mutex `json:"-"` // Mutex for concurrent access control; excluded from JSON output.
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
