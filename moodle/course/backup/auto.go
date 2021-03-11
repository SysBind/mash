// moodle/course/backup/auto.go

package backup

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/sysbind/mash/moodle/config"
	"github.com/sysbind/mash/moodle/database"
)

type Storage int

const (
	STORAGE_COURSE               Storage = 0
	STORAGE_DIRECTORY                    = 1
	STORAGE_COURSE_AND_DIRECTORY         = 2
)

type AutoBackup struct {
	active  bool
	maxkept int
	storage Storage
	cfg     config.Config
}

func LoadAutoBackup(cfg config.Config) (ab AutoBackup, err error) {
	ab.cfg = cfg

	intVal, err := strconv.Atoi(cfg.GetPluginConf("backup", "backup_auto_active"))
	if err != nil {
		return
	}
	ab.active = intVal > 0

	ab.maxkept, err = strconv.Atoi(cfg.GetPluginConf("backup", "backup_auto_max_kept"))
	if err != nil {
		return
	}

	intVal, err = strconv.Atoi(cfg.GetPluginConf("backup", "backup_auto_storage"))
	if err != nil {
		return
	}
	ab.storage = Storage(intVal)
	return
}

// PreFlight checks some basic settings to see if autobackup should run
func (ab AutoBackup) PreFlight() (err error) {
	if !ab.active {
		err = errors.New("Automated Course Backups not enabled")
	}
	return
}

// Run starts the automated backup process
func (ab AutoBackup) Run() (err error) {
	var ids []int

	if ids, err = ab.getCourses(); err != nil {
		log.Fatal(err)
	}

	for i := range ids {
		ab.backupCourse(ids[i])
	}

	return
}

// getCourses returns course ids to backup
func (ab AutoBackup) getCourses() (ids []int, err error) {
	query := fmt.Sprintf("SELECT id FROM mdl_course ORDER BY id DESC")
	var db database.Database = ab.cfg.DB()

	rows, err := db.Query(query)
	if err != nil {
		return
	}

	defer rows.Close()

	ids = make([]int, 0)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		ids = append(ids, id)
	}
	// Check for errors from iterating over rows.
	err = rows.Err()

	return
}

// backupCourse creates single course backup.
func (ab AutoBackup) backupCourse(id int) (err error) {
	fmt.Printf("backing up %d \n", id)

	cmd := exec.Command("php", "admin/cli/automated_backup_single.php", strconv.Itoa(id))
	out, err := cmd.CombinedOutput()
	if err != nil {
		if out != nil {
			fmt.Printf("combined out:\n%s\n", string(out))
		}
		fmt.Printf("cmd.Run() failed with %v\n", err)
		fmt.Printf("backup of %d failed ! \n", id)
	}

	fmt.Printf("backup of %d finished \n", id)

	ab.removeExcessBackups(id)

	return
}

// removeExcessBackups deletes old backups according to auto backup settings
// logic copied from backup/util/helper/backup_cron_helper.class.php::remove_excess_backups
func (ab AutoBackup) removeExcessBackups(id int) (err error) {
	if ab.maxkept == 0 {
		return
	}

	if ab.storage == STORAGE_COURSE || ab.storage == STORAGE_COURSE_AND_DIRECTORY {
		err = ab.removeExcessBackupsFromCourse(id)
	}

	if ab.storage == STORAGE_DIRECTORY || ab.storage == STORAGE_COURSE_AND_DIRECTORY {
		err = ab.removeExcessBackupsFromDir(id)
	}
	return
}

// removeExcessBackupsFromCourse removes old backups from course stroage area
func (ab AutoBackup) removeExcessBackupsFromCourse(id int) (err error) {
	return
}

// removeExcessBackupsFromDir removes old backups from backup dir
func (ab AutoBackup) removeExcessBackupsFromDir(id int) (err error) {
	return
}
