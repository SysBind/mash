package main

import (
	"fmt"
	"log"

	"github.com/sysbind/moodle-automated-course-backup/moodle"
)

func main() {
	cfg, err := moodle.Parse("config.php")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)
}
