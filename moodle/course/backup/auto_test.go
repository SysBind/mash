// moodle/course/backup/auto_test.go

package backup

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/sysbind/moodle-automated-course-backup/moodle/config"
	"github.com/sysbind/moodle-automated-course-backup/moodle/database"
)

func TestPreFlight(t *testing.T) {
	is := is.New(t)

	cfg, err := config.Parse("../../testdata/config.php")
	db, err := database.Open(context.Background(), cfg.DriverName(), cfg.DSN())

	is.NoErr(err)
	cfg.SetDatabase(db)

	cfg.SetPluginConf("backup", "backup_auto_active", "0")

	err = PreFlight(cfg)

	is.True(err != nil) // PreFlight should fail when backup_auto_active is 0

	return
}
