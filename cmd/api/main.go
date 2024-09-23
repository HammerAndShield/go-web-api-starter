package main

import (
	"context"
	"fmt"
	"go-web-api-starter/internal/apiutils"
	"os"
)

type application struct {
	config *apiutils.ApiConfig
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Getenv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
