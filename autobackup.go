package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/sysbind/moodle-automated-course-backup/moodle/config"
	"github.com/sysbind/moodle-automated-course-backup/moodle/database"
)

func main() {
	cfg, err := config.Parse("config.php")

	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	db := database.Open(cfg.DriverName(), cfg.DSN(), ctx)
	defer db.Close()

	appSignal := make(chan os.Signal, 3)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		select {
		case <-appSignal:
			stop()
		}
	}()
}
