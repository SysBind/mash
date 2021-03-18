// moodle/course/backup/backup.go

package backup

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sysbind/mash/moodle/database"
)

type Status uint // Backup Status

const (
	STATUS_OK         Status = 1 // Course automated backup completed successfully
	STATUS_ERROR      Status = 0 // Course automated backup errored
	STATUS_UNFINISHED Status = 2 // Course automated backup never finished
	STATUS_SKIPPED    Status = 3 // Course automated backup was skipped
	STATUS_WARNING    Status = 4 // Course automated backup had warnings
	STATUS_NOTYETRUN  Status = 5 // Course automated backup has yet to be run
)

// CourseBackupRec stores state of single course backup run (mdl_backup_courses)
type CourseBackupRec struct {
	Id        uint64
	CourseId  uint64
	StartTime uint64
	EndTime   uint64
	Status    Status
	Message   sql.NullString
}

// StartBackupRec records start time of backup and returns CourseBackupRec
//                for further updates
func startBackupRec(db database.Database, cid uint64) (cb CourseBackupRec, err error) {
	cb, err = getBackupRec(db, cid)
	if err != nil {
		return
	}
	cb.StartTime = uint64(time.Now().Unix())
	cb.EndTime = 0
	cb.Status = STATUS_UNFINISHED
	err = cb.updateRow(db)

	return
}

// getBackupRec gets record of course backup from mdl_backup_courses
//              will create the record if not exist yet
func getBackupRec(db database.Database, cid uint64) (cb CourseBackupRec, err error) {
	cb.CourseId = cid
	query := fmt.Sprintf("SELECT id,laststarttime,lastendtime,laststatus,message FROM mdl_backup_courses WHERE courseid=%d", cid)

	err = db.QueryRow(query).Scan(&cb.Id,
		&cb.StartTime,
		&cb.EndTime,
		&cb.Status,
		&cb.Message)

	if err != nil {
		if err == sql.ErrNoRows {
			err = cb.insertRow(db)
		}
		return
	}

	return
}

func (cb *CourseBackupRec) insertRow(db database.Database) (err error) {
	query := fmt.Sprintf("INSERT INTO mdl_backup_courses(courseid) VALUES(%d)", cb.CourseId)
	var result sql.Result

	result, err = db.Exec(query)
	if err != nil {
		return
	}

	var rows int64
	rows, err = result.RowsAffected()
	if err != nil {
		return
	}

	if rows != 1 {
		err = fmt.Errorf("course/backup: insertRow() -  RowsAffected is %d", rows)
		return
	}

	lastid, err := result.LastInsertId()
	cb.Id = uint64(lastid)

	return
}

func (cb *CourseBackupRec) updateRow(db database.Database) (err error) {
	message := ""
	if cb.Message.Valid {
		message = cb.Message.String
	}
	query := fmt.Sprintf("UPDATE mdl_backup_courses SET laststarttime=%d, lastendtime=%d, laststatus=%d, message='%s' WHERE id=%d", cb.StartTime, cb.EndTime, cb.Status, message, cb.Id)
	fmt.Println("updateRow: ", query)

	var result sql.Result
	result, err = db.Exec(query)
	if err != nil {
		return
	}

	var rows int64
	rows, err = result.RowsAffected()
	if err != nil {
		return
	}

	if rows != 1 {
		err = fmt.Errorf("course/backup: updateRow() - RowsAffected is %d", rows)
	}

	return
}
