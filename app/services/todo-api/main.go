package main

import (
	"errors"
	"expvar"
	"fmt"
	"github.com/ardanlabs/conf/v3"
	"log"
	"os"
	"runtime"
	"time"
	"todo-svc/app/services/todo-api/handlers"
	"todo-svc/business/database"
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

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
		}
		Auth struct {
			KeysFolder string `conf:"default:zarf/keys/"`
		}
		DB struct {
			User         string `conf:"default:postgres"`
			Password     string `conf:"default:postgres,mask"`
			Host         string `conf:"default:localhost"`
			Name         string `conf:"default:postgres"`
			MaxIdleConns int    `conf:"default:0"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}
	}{
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

	logg.Println("Starting todo-api service")
	defer logg.Println("Shutdown complete")

	expvar.NewString("build").Set(build)

	// Setup persistence
	pers, err := database.NewMemoryPersistence()
	if err != nil {
		return fmt.Errorf("error initializing database connection: %w", err)
	}

	// Should setup server mux
	if err := handlers.SetupHandlers(logg, pers); err != nil {
		return fmt.Errorf("error initializing handlers: %w", err)
	}

	return nil
}
