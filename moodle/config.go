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
	switch field {
	case "dbtype":
		cfg.dbtype = value
	case "dbhost":
		cfg.dbhost = value
	case "dbname":
		cfg.dbname = value
	case "dbuser":
		cfg.dbuser = value
	case "dbpass":
		cfg.dbpass = value
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
			cfg.assignFieldValue(parseLine(line))
		}
	}
	err = scanner.Err()
	return
}

// Parses one line of config.php into field and value
func parseLine(line string) (field, value string) {
	field_n_value := strings.Split(line, "->")[1]
	field = strings.TrimSpace(strings.Split(field_n_value, "=")[0])
	value = cleanValue(strings.TrimSpace(strings.Split(field_n_value, "=")[1]))
	return
}

// Clean a value from trailing ";", possible comments, and surrounding "'"
func cleanValue(value string) (retval string) {
	retval = strings.Split(value, ";")[0]
	retval = strings.Trim(retval, "'")
	return
}
