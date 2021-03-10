package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/sysbind/moodle-automated-course-backup/moodle/config"
	"github.com/sysbind/moodle-automated-course-backup/moodle/course/backup"
	"github.com/sysbind/moodle-automated-course-backup/moodle/database"
)

func main() {
	cfg, db, stop := setup()
	defer stop()
	defer db.Close()

	err := backup.PreFlight(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = backup.Run(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

// setup parses config, check db connection and creates context
func setup() (cfg config.Config, db database.Database, stop context.CancelFunc) {
	cfg, err := config.Parse("config.php")

	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := context.WithCancel(context.Background())

	db, err = database.Open(ctx, cfg.DriverName(), cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}

	cfg.SetDatabase(db)

	// Catch OS interupt to cancel the context (stop())
	appSignal := make(chan os.Signal, 3)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		select {
		case <-appSignal:
			stop()
		}
	}()

	return
}
