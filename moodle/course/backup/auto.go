// moodle/course/backup/auto.go

package backup

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/sysbind/mash/moodle"
	"github.com/sysbind/mash/moodle/config"
	"github.com/sysbind/mash/moodle/database"
	"github.com/sysbind/mash/moodle/storage"
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
	dest    string // dest if storage is *_DIRECTORY
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
	ab.dest = cfg.GetPluginConf("backup", "backup_auto_destination")
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
		return
	}

	for i := range ids {
		err = ab.backupCourse(ids[i])
		if err != nil {
			return
		}
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

	if err = cmd.Start(); err != nil {
		return
	}

	if err = cmd.Wait(); err != nil {
		fmt.Printf("backup of %d failed !! \n", id)
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
			}
		}
		// Only way to extract the actual stderr/stdout -
		//   run it again with CombinedOutput
		cmd := exec.Command("php", "admin/cli/automated_backup_single.php", strconv.Itoa(id))
		out, _ := cmd.CombinedOutput()
		log.Println(string(out))

		return
	}

	fmt.Printf("backup of %d finished \n", id)

	err = ab.removeExcessBackups(id)

	return
}

// removeExcessBackups deletes old backups according to auto backup settings
// logic copied from backup/util/helper/backup_cron_helper.class.php::remove_excess_backups
func (ab AutoBackup) removeExcessBackups(id int) (err error) {
	if ab.maxkept == 0 {
		return
	}

	if ab.storage == STORAGE_COURSE || ab.storage == STORAGE_COURSE_AND_DIRECTORY {
		if err = ab.removeExcessBackupsFromCourse(id); err != nil {
			return
		}
	}

	if ab.storage == STORAGE_DIRECTORY || ab.storage == STORAGE_COURSE_AND_DIRECTORY {
		err = ab.removeExcessBackupsFromDir(id)
	}
	return
}

// removeExcessBackupsFromCourse removes old backups from course stroage area
func (ab AutoBackup) removeExcessBackupsFromCourse(id int) (err error) {
	var files []storage.StoredFile

	if files, err = ab.getAutoBackupsFromCourse(id); err != nil {
		return
	}

	// drop last maxkept elements
	files = files[:len(files)-ab.maxkept]

	for _, file := range files {
		if err = file.Delete(ab.cfg); err != nil {
			return
		}
	}

	return
}

// removeExcessBackupsFromDir removes old backups from backup dir
func (ab AutoBackup) removeExcessBackupsFromDir(id int) (err error) {
	glob := fmt.Sprintf("%s/backup-moodle2-course-%d-*.mbz",
		ab.dest,
		id)
	var files []string
	files, err = filepath.Glob(glob)

	if err != nil {
		return
	}

	// drop last maxkept elements
	files = files[:len(files)-ab.maxkept]

	for _, file := range files {
		err = os.Remove(file)
		fmt.Println("removeExcessBackupsFromDir", file)
	}
	return
}

func (ab AutoBackup) getAutoBackupsFromCourse(id int) (files []storage.StoredFile, err error) {
	var db database.Database = ab.cfg.DB()
	var cctx moodle.Context

	cctx, err = moodle.CourseContext(db, id)
	if err != nil {
		return
	}

	query := fmt.Sprintf("SELECT id, filename, contenthash, contextid, component, filearea, timecreated FROM mdl_files WHERE contextid=%d AND component='%s' AND filearea='%s' AND NOT filename = '.' ORDER BY timecreated ASC", cctx.Id, "backup", "automated")

	rows, err := db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var file storage.StoredFile
		if err = rows.Scan(&file.Id,
			&file.FileName,
			&file.ContentHash,
			&file.ContextId,
			&file.Component,
			&file.FileArea,
			&file.TimeCreated); err != nil {
			return
		}
		files = append(files, file)
	}
	// Check for errors from iterating over rows.
	err = rows.Err()

	return
}
