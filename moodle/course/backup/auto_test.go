// moodle/course/backup/auto_test.go

package backup

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/sysbind/mash/moodle/config"
	"github.com/sysbind/mash/moodle/database"
)

func TestPreFlight(t *testing.T) {
	is := is.New(t)

	cfg, err := config.Parse("../../testdata/config.php")
	is.NoErr(err) // parse config

	db, err := database.Open(context.Background(), cfg.DriverName(), cfg.DSN())
	is.NoErr(err) // open database

	cfg.SetDatabase(db)

	cfg.SetPluginConf("backup", "backup_auto_active", "0")
	ab, err := LoadAutoBackup(cfg)
	is.NoErr(err) // can LoadAutoBackup settings

	err = ab.PreFlight()
	is.True(err != nil) // PreFlight should fail when backup_auto_active is 0

	return
}

func TestRun(t *testing.T) {
	is := is.New(t)

	cfg, err := config.Parse("../../testdata/config.php")

	is.True(cfg.DirRoot() != "") // dirroot must be set in config.php for unit testing

	db, err := database.Open(context.Background(), cfg.DriverName(), cfg.DSN())

	is.NoErr(err)
	cfg.SetDatabase(db)

	dest, err := createTempBackupDir()
	fmt.Printf("Created Temp Backup Dir %s", dest)
	is.NoErr(err)
	defer removeTempBackupDir(dest)

	cfg.SetPluginConf("backup", "backup_auto_active", "1")
	cfg.SetPluginConf("backup", "backup_auto_destination", dest)
	cfg.SetPluginConf("backup", "backup_auto_storage",
		fmt.Sprintf("%d", STORAGE_COURSE_AND_DIRECTORY))

	ab, err := LoadAutoBackup(cfg)
	is.NoErr(err) // can LoadAutoBackup settings

	err = os.Chdir(cfg.DirRoot()) //
	is.NoErr(err)                 // change dir to moodle root should not error
	err = exec.Command("php", "admin/cli/purge_caches.php").Run()
	is.NoErr(err) // purge caches should not error

	err = ab.PreFlight()
	is.NoErr(err) // PreFlight should not return error

	err = ab.Run()
	is.NoErr(err) // Run should not return error

	files, _ := ioutil.ReadDir(dest)
	is.True(len(files) == 1) // One backup file should have been created

	// run again (sleep 1 minute)
	time.Sleep(time.Minute)
	err = ab.Run()
	is.NoErr(err) // Second Run should not return error

	files, _ = ioutil.ReadDir(dest)
	is.True(len(files) == 1) // One backup file should exist after second Run (last one should have been deleted)

	return
}

func createTempBackupDir() (dir string, err error) {
	dir, err = os.Getwd()
	if err != nil {
		return
	}
	dir = fmt.Sprintf("%s/%s", dir, "tmp")
	os.Mkdir(dir, os.ModePerm)
	return
}

func removeTempBackupDir(dir string) {
	os.RemoveAll(dir)
}
