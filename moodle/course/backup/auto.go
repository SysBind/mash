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
	"time"

	"github.com/sysbind/mash/moodle"
	"github.com/sysbind/mash/moodle/config"
	"github.com/sysbind/mash/moodle/course"
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
	active        bool
	maxkept       int
	storage       Storage
	dest          string // dest if storage is *_DIRECTORY
	skipmodifprev bool   // skip if not modified since last backup
	cfg           config.Config
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
	var ids []uint64

	if ids, err = ab.getCourses(); err != nil {
		return
	}

	courses := make(chan uint64, len(ids))
	results := make(chan int, len(ids))

	// create workers
	for w := 1; w <= 3; w++ {
		go ab.worker(w, courses, results)
	}

	// feed courses into channel
	for _, course := range ids {
		courses <- course
	}
	close(courses)

	// wait for all backups to complete
	for a := 1; a <= len(ids); a++ {
		<-results
	}

	log.Println("Auto Backup Done")

	return
}

func (ab AutoBackup) worker(id int, courses <-chan uint64, results chan<- int) {
	for cid := range courses {
		err := ab.backupCourse(cid)
		if err != nil {
			log.Println("error on course", cid, err)
			results <- 2
			continue
		}
		results <- 0
	}
}

// getCourses returns course ids to backup
func (ab AutoBackup) getCourses() (ids []uint64, err error) {
	query := fmt.Sprintf("SELECT id FROM mdl_course ORDER BY id DESC")
	var db database.Database = ab.cfg.DB()

	rows, err := db.Query(query)
	if err != nil {
		return
	}

	defer rows.Close()

	ids = make([]uint64, 0)
	for rows.Next() {
		var id uint64
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
func (ab AutoBackup) backupCourse(id uint64) (err error) {
	// mdl_backup_courses book keeping:
	var db database.Database = ab.cfg.DB()
	var backupRec CourseBackupRec
	backupRec, err = getBackupRec(db, id)
	if err != nil {
		return
	}
	defer backupRec.updateRow(db)

	// Skip unmodified since last backup?
	if ab.skipmodifprev && !course.ModifiedSince(db, id, backupRec.EndTime) {
		backupRec.Status = STATUS_SKIPPED
		backupRec.Message.String = fmt.Sprintf("Not modified since last backup in %s", time.Unix(int64(backupRec.EndTime), 0))
		return
	}
	backupRec.StartTime = uint64(time.Now().Unix())
	backupRec.EndTime = 0
	backupRec.Status = STATUS_UNFINISHED
	backupRec.Message.String = "Backup Start"

	cmd := exec.Command("php", "admin/cli/automated_backup_single.php", strconv.FormatUint(id, 10))
	if err = cmd.Start(); err != nil {
		return
	}

	if err = cmd.Wait(); err != nil {
		fmt.Printf("backup of %d failed !! \n", id)
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
			}
		}
		// Only way to extract the actual stderr/stdout -
		//   run it again with CombinedOutput
		cmd := exec.Command("php", "admin/cli/automated_backup_single.php", strconv.FormatUint(id, 10))
		out, _ := cmd.CombinedOutput()
		log.Println(string(out))

		backupRec.Status = STATUS_ERROR
		backupRec.Message.String = string(out)
		backupRec.EndTime = uint64(time.Now().Unix())
		return
	}
	err = ab.removeExcessBackups(id)
	backupRec.EndTime = uint64(time.Now().Unix())
	backupRec.Message.String = ""
	backupRec.Status = STATUS_OK
	return
}

// removeExcessBackups deletes old backups according to auto backup settings
// logic copied from backup/util/helper/backup_cron_helper.class.php::remove_excess_backups
func (ab AutoBackup) removeExcessBackups(id uint64) (err error) {
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
func (ab AutoBackup) removeExcessBackupsFromCourse(id uint64) (err error) {
	var files []storage.StoredFile

	if files, err = ab.getAutoBackupsFromCourse(id); err != nil {
		return
	}

	if len(files) <= ab.maxkept {
		return // nothing to delete
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
func (ab AutoBackup) removeExcessBackupsFromDir(id uint64) (err error) {
	glob := fmt.Sprintf("%s/backup-moodle2-course-%d-*.mbz",
		ab.dest,
		id)
	var files []string
	files, err = filepath.Glob(glob)

	if err != nil {
		return
	}

	if len(files) <= ab.maxkept {
		return // nothing to delete
	}

	// drop last maxkept elements
	files = files[:len(files)-ab.maxkept]

	for _, file := range files {
		err = os.Remove(file)
	}
	return
}

func (ab AutoBackup) getAutoBackupsFromCourse(id uint64) (files []storage.StoredFile, err error) {
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

// LoadAutoBackup Loads auto backup settings from database
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

	intVal, err = strconv.Atoi(cfg.GetPluginConf("backup", "backup_auto_skip_modif_prev"))
	if err != nil {
		return
	}
	ab.skipmodifprev = intVal > 0

	return
}
