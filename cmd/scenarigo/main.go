package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/bilus/scenarigo/cmd/scenarigo/cmd"
)

func main() {
	if err := run(); err != nil {
		if errors.Is(err, cmd.ErrTestFailed) {
			os.Exit(10)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	return cmd.Execute(ctx)
}
