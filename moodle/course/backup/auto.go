// moodle/course/backup/auto.go

package backup

import (
	"errors"
	"strconv"

	"github.com/sysbind/moodle-automated-course-backup/moodle/config"
)

// PreFlight checks some basic settings to see if autobackup should run
func PreFlight(cfg config.Config) (err error) {
	isActive := cfg.GetPluginConf("backup", "backup_auto_active")

	var active int
	if active, err = strconv.Atoi(isActive); err == nil && active == 0 {
		err = errors.New("Automated Course Backups not enabled")
	}

	return
}
