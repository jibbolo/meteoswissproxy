package main

import (
	"cmp"
	"log/slog"
	"os"
	"time"
)

const readTimeout = time.Second * 5

func main() {
	port := cmp.Or(os.Getenv("PORT"), "8080")
	if err := startHttpServer(port); err != nil {
		slog.Error("http server failed", "error", err)
		os.Exit(1)
	}
}
