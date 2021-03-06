package main

import (
	"errors"
	"expvar"
	"fmt"
	"github.com/ardanlabs/conf/v3"
	"log"
	"os"
	"runtime"
	"todo-svc/app/services/todo-api/config"
	"todo-svc/app/services/todo-api/handlers"
	"todo-svc/business/database"
	"todo-svc/foundation/management"
)

var build = "develop"

func main() {
	logg := log.New(os.Stderr, "[todo-api] ", log.LstdFlags)

	if err := run(logg); err != nil {
		logg.Fatalf("todo-api failed to start: %v\n", err)
	}
}

func run(logg *log.Logger) error {
	cpus := runtime.NumCPU()
	logg.Printf("Have %d CPUs", cpus)

	cfg := config.Config{
		Version: conf.Version{
			Build: build,
			Desc:  "Todo Item API",
		},
	}

	const prefix = "TODO"
	if help, err := conf.Parse(prefix, &cfg); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	logg.Println(out)
	expvar.NewString("build").Set(build)

	// Setup persistence
	pers, err := database.NewMemoryPersistence()
	if err != nil {
		return fmt.Errorf("error initializing database connection: %w", err)
	}

	// Setup server mux
	app, err := handlers.SetupHandlers(logg, cfg, pers)
	if err != nil {
		return fmt.Errorf("error initializing handlers: %w", err)
	}

	// Setup context cancellation and start the server
	ctx := management.InterruptContext(logg)
	httpDone := management.AsyncListen(ctx, cfg, logg, app)

	// Wait for the HTTP side to complete shutdown
	<-httpDone
	return nil
}
