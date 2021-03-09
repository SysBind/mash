// moodle/config/plugins.go
package config

import (
	"fmt"
	"log"
)

func (cfg Config) GetPluginConf(plugin, name string) (value string) {
	query := fmt.Sprintf("SELECT value FROM mdl_config_plugins WHERE plugin='%s' AND name='%s'", plugin, name)

	err := cfg.db.QueryRow(query).Scan(&value)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (cfg Config) SetPluginConf(plugin, name, value string) {
}
