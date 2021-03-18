// moodle/course/backup/backup.go

package backup

import (
	"database/sql"
	"fmt"

	"github.com/sysbind/mash/moodle/database"
)

type Status int // Backup Status

const (
	STATUS_OK         Status = 1 // Course automated backup completed successfully
	STATUS_ERROR      Status = 0 // Course automated backup errored
	STATUS_UNFINISHED Status = 2 // Course automated backup never finished
	STATUS_SKIPPED    Status = 3 // Course automated backup was skipped
	STATUS_WARNING    Status = 4 // Course automated backup had warnings
	STATUS_NOTYETRUN  Status = 5 // Course automated backup has yet to be run
)

type CourseBackup struct {
	Id        int64
	CourseId  int64
	StartTime int
	EndTime   int
	Status    int
	Message   string
}

func GetBackup(db database.Database, cid int64) (cb CourseBackup, err error) {
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

func (cb *CourseBackup) insertRow(db database.Database) (err error) {
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

	cb.Id, err = result.LastInsertId()

	return
}

func (cb *CourseBackup) updateRow(db database.Database) (err error) {
	query := fmt.Sprintf("UPDATE mdl_backup_courses SET StartTime=%d, EndTime=%d, Status=%d, Message='%s' WHERE id=%d", cb.StartTime, cb.EndTime, cb.Status, cb.Message, cb.Id)

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
