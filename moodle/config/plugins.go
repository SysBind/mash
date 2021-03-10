// moodle/config/plugins.go
package config

import (
	"fmt"
	"log"
)

// GetPluginConf gets single configuration value from config_plugins table
func (cfg Config) GetPluginConf(plugin, name string) (value string) {
	query := fmt.Sprintf("SELECT value FROM mdl_config_plugins WHERE plugin='%s' AND name='%s'", plugin, name)

	err := cfg.db.QueryRow(query).Scan(&value)
	if err != nil {
		log.Fatalf("GetPluginConf: %v (%s)", err, query)
	}
	return
}

// SetPluginConf sets single configuration value in config_plugins table
func (cfg Config) SetPluginConf(plugin, name, value string) {
	// get current config
	curval := cfg.GetPluginConf(plugin, name)
	if value == curval {
		// no modification needed
		return
	}

	query := fmt.Sprintf("UPDATE mdl_config_plugins SET value='%s' WHERE plugin='%s' AND name='%s'", value, plugin, name)

	fmt.Printf("Executing %s\n", query)

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
	return
}
