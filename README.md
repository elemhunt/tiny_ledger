
# Tiny Ledger API

A lightweight RESTful API service for recording and viewing basic financial transactions (deposits and withdrawals), checking balance, and reviewing transaction history.

This project is a take-home assignment designed to demonstrate practical API design, concurrency safety, and clean code structure using idiomatic Go.

---

## Features

- Record money movements (deposits and withdrawals)
- View current balance
- View full transaction history
- Thread-safe, in-memory data storage
- Simple and fast local deployment

---

## Project Structure

```

.
├── cmd/api             # Entry point for the server
├── config              # (Future) configuration setup
├── internal
│   ├── handler         # HTTP handlers, ledger business logic
│   └── server          # Router and server setup
└── tests               # Unit and integration tests

````

---

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git

---

### Environment Setup

Before running or building the project, configure your environment:

1. Copy the environment template:

    ```bash
    cp env.tpl env.dev
    ```

2. Create a symbolic link for `.env`:

    ```bash
    ln -s env.dev .env
    ```

3. Make changes only to `env.dev`. Both `.env` and `env.dev` are excluded from version control.

4. Set the desired port in `env.dev`:

    ```env
    PORT=8080
    ```

---

### Run Locally

```bash
# Clone the repository
git clone https://github.com/elemhunt/tiny_ledger.git
cd tiny_ledger

# Run the API
go run ./cmd/api
````

By default, the server will start on `http://localhost:8080`.

---

### Building and Running in Production

To build a statically compiled binary suitable for production use:

```bash
go build -o tiny-ledger ./cmd/api
```

This creates an executable named `tiny-ledger` in the current directory.

You can then run it:

```bash
./tiny-ledger
```

Make sure `.env` is configured and present in the same directory, or exported environment variables are available. Logs and errors will be printed to standard output.

---

## API Endpoints

### POST /ledger/transactions (Deposit)

Records a new transaction.

**Request Body:**

```json
{
  "type": "deposit", 
  "amount": 100.0
}
```

**Response (201 Created):**

```json
{
  "id": 1,
  "type": "deposit",
  "amount": 100,
  "timestamp": "2025-05-11T12:00:00Z"
}
```

---

### Example: POST /ledger/transactions (Withdrawal)

For a **withdrawal** transaction:

**Request Body:**

```json
{
  "type": "withdrawal",  // "withdrawal" type
  "amount": 50.0         // Withdrawal amount
}
```

**Response (201 Created):**

```json
{
  "id": 2,
  "type": "withdrawal",
  "amount": 50,
  "timestamp": "2025-05-11T12:15:00Z"
}
```

---

### GET /ledger/balance

Returns the current account balance.

**Response:**

```json
{
  "balance": 100.0,
  "checked_at": "2025-05-11T12:01:00Z"
}
```

---

### GET /ledger/transaction\_history

Returns a list of all transactions.

**Response:**

```json
[
  {
    "id": 1,
    "type": "deposit",
    "amount": 100,
    "timestamp": "2025-05-11T12:00:00Z"
  },
  {
    "id": 2,
    "type": "withdrawal",
    "amount": 50,
    "timestamp": "2025-05-11T12:15:00Z"
  }
]
```

---

## Running Tests

```bash
go test ./tests/...
```

Tests cover route functionality, expected response codes, and core behaviors.

---

## Framework Choice

We use [`chi`](https://github.com/go-chi/chi) for routing due to its:

* Minimal and idiomatic API design
* Modular middleware support (logging, recovery, CORS, etc.)
* Robust routing features with nested routes
* Clean integration with Go’s standard library

---

## Future Improvements

* Persistent storage using SQLite or Postgres
* Token-based authentication
* Filtering and pagination on transaction history
* Docker support
* Structured logging and observability tooling

---

## Assumptions

* The application is single-user (multi-user support is not implemented)
* Transactions are guarded by mutexes to ensure thread-safety
* Negative amounts and invalid transaction types are rejected
* Balances are floats (no high-precision handling)

---

## Tech Stack

* Go
* `chi` for routing
* `testify` and standard `testing` packages for tests

---

## Contact

For questions or suggestions, please open an issue on GitHub or contact [@elemhunt](https://github.com/elemhunt).

---

## License

This project is licensed under the MIT License.



