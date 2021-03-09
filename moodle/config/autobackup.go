// moodle/config
package config

import (
	"fmt"
	"log"
	"strconv"

	"github.com/sysbind/moodle-automated-course-backup/moodle/database"
)

// config.AutoBackup struct holds information from
// moodle's config_plugins table starting with backup_auto
type AutoBackup struct {
	Active bool
}

func (settings *AutoBackup) assignFieldValue(field, value string) {
	switch field {
	case "backup_auto_active":
		if intVal, err := strconv.Atoi(value); err == nil && intVal > 0 {
			settings.Active = true
		} else {
			settings.Active = false
		}

	}
}

func (settings *AutoBackup) String() string {
	return fmt.Sprintf("Active: %t", settings.Active)
}

func GetAutoBackupSettings(db database.Database) (settings AutoBackup) {
	rows, err := db.Query("SELECT name, value FROM mdl_config_plugins WHERE plugin='backup' AND name LIKE 'backup_auto%'")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name, value string
		if err := rows.Scan(&name, &value); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		settings.assignFieldValue(name, value)
	}
	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return settings
}
