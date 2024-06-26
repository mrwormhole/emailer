package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/mrwormhole/emailer/brevo"
)

var debugEnabled = flag.Bool("debug", false, "in debug environment")

const (
	defaultPort = "5555"
	// Providers Listed below
	providerBrevo = "brevo"
)

func main() {
	flag.Parse()

	var logHandler slog.Handler
	if *debugEnabled {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	} else {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}
	slog.SetDefault(slog.New(logHandler))

	key, ok := os.LookupEnv("API_KEY")
	if !ok {
		slog.Error("API_KEY not found in env")
		os.Exit(1)
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}
	portNum, err := strconv.Atoi(port)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, fmt.Sprintf("strconv.Atoi(%q)", port), slog.String("err", err.Error()))
		os.Exit(1)
	}
	provider, ok := os.LookupEnv("PROVIDER")
	if !ok {
		provider = providerBrevo
	}

	c := retryablehttp.NewClient()
	c.RetryMax = 3
	httpClient := c.StandardClient()
	httpClient.Timeout = 10 * time.Second

	var handler http.HandlerFunc
	switch {
	case strings.EqualFold(provider, providerBrevo):
		sender, err := brevo.New(key, httpClient)
		if err != nil {
			slog.LogAttrs(context.Background(), slog.LevelError, "brevo.New()", slog.String("err", err.Error()))
		}
		handler = brevo.EmailHandler(sender)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /email", handler)

	srv := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", portNum),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		slog.Debug(fmt.Sprintf("server started at localhost:%d", portNum))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.LogAttrs(context.Background(), slog.LevelError, "srv.ListenAndServe()", slog.String("err", err.Error()))
		}
	}()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, os.Interrupt)
	<-wait

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "srv.Shutdown()", slog.String("err", err.Error()))
	}
}
