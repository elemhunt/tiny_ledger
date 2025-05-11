package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/elemhunt/tiny_ledger/internal/server"
)

func main() {

	server := server.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := server.Start(ctx)
	if err != nil {
		fmt.Println("failed server init:", err)
	}

}
