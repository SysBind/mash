// moodle/course/eventlog.go

package course

import (
	"fmt"

	"github.com/sysbind/mash/moodle/database"
)

func ModifiedSince(db database.Database, cid uint64, since uint64) bool {

	query := fmt.Sprintf("SELECT COUNT(id) FROM mdl_logstore_standard_log WHERE courseid = %d AND timecreated > %d AND NOT crud = 'r' AND NOT target = 'course_backup'", cid, since)
	var count int
	_ = db.QueryRow(query).Scan(&count)
	if count > 0 {
		return true
	}
	query = fmt.Sprintf("SELECT COUNT(id) FROM mdl_logstore_standard_log WHERE courseid = %d AND target = 'course_backup'", cid)
	_ = db.QueryRow(query).Scan(&count)
	return count == 0

}
