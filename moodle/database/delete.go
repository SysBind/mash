// moodle/database/delete.go

package database

import (
	"fmt"
	"log"
)

func (db Database) DeleteRecord(table string, id int) {
	query := fmt.Sprintf("DELETE FROM mdl_%s WHERE id=%d", value, plugin, name)

	result, err := cfg.db.Exec(query)
	if err != nil {
		log.Fatalf("SetPluginConf: %v (%s)", err, query)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("SetPluginConf: %v (%s)", err, query)
	}

	if rows != 1 {
		log.Fatalf("SetPluginConf: expected single row affected, got %d rows affected (%s)", rows, query)
	}
}
