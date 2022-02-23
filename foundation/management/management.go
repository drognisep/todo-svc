// Package management provides some useful utility functionality that can be used to express common server management patterns.
package management

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"todo-svc/app/services/todo-api/config"
)

// InterruptContext will setup context cancellation in the event that a SIGINT or SIGTERM signal is received.
func InterruptContext(logg *log.Logger) context.Context {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		hasCancelled := false
		for {
			select {
			case <-ch:
				if !hasCancelled {
					logg.Println("Received signal, initiating controlled stop")
					cancel()
					hasCancelled = true
					continue
				}
				logg.Fatalln("Received second signal, stopping now")
			}
		}
	}()

	return ctx
}

func AsyncListen(ctx context.Context, cfg config.Config, logg *log.Logger, app *fiber.App) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer logg.Println("Shutdown complete")
		if cfg.Auth.Mode == "dev" {
			logg.Println("Starting server in DEV MODE...")
			if err := app.Listen(cfg.Web.APIHost); err != nil {
				logg.Printf("Error starting server: %v\n", err)
			}
		} else {
			logg.Println("Starting server...")
			if err := app.ListenTLS(cfg.Web.APIHost, cfg.Web.CertFile, cfg.Web.KeyFile); err != nil {
				logg.Printf("Error starting server: %v\n", err)
			}
		}
		close(ch)
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				logg.Println("Received controlled stop signal")
				if err := app.Shutdown(); err != nil {
					logg.Printf("Error shutting down server: %v\n", err)
				}
			}
		}
	}()
	return ch
}
