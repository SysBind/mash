package moodle

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config struct holds information from moodle's config.php
type Config struct {
	dbtype string
	dbhost string
	dbname string
	dbuser string
	dbpass string
}

// String representation of Config struct
func (cfg Config) String() string {
	return fmt.Sprintf("%s://%s:%s@%s/%s",
		cfg.dbtype,
		cfg.dbuser,
		cfg.dbpass,
		cfg.dbhost,
		cfg.dbname)
}

func (cfg *Config) assignFieldValue(field, value string) {
	fmt.Printf("Assign %s = %s\n", field, value)
	switch field {
	case "dbtype":
		fmt.Println("Assign dbtype " + value)
		cfg.dbtype = value
	case "dbhost":
		cfg.dbhost = value
	case "dbname":
		cfg.dbname = value
	case "dbuser":
		cfg.dbuser = value
	case "dbpass":
		cfg.dbpass = value
	default:
		fmt.Printf("Unknown config %s", value)
	}
}

// Parse reads config.php and decodes it into Config struct
func Parse(filename string) (cfg Config, err error) {
	f, err := os.Open(filename)

	if err != nil {
		return
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "$CFG->") {
			field := strings.Split(line, "->")[1]
			value := strings.Split(field, "=")[1]
			cfg.assignFieldValue(field, value)
		}
	}
	err = scanner.Err()
	return
}
