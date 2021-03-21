// moodle/context.go

package moodle

import (
	"database/sql"
	"fmt"

	"github.com/sysbind/mash/moodle/database"
)

type ContextLevel int

const (
	CONTEXT_SYSTEM    ContextLevel = 10
	CONTEXT_USER                   = 30
	CONTEXT_COURSECAT              = 40
	CONTEXT_COURSE                 = 50
	CONTEXT_MODULE                 = 70
	CONTEXT_BLOCK                  = 80
)

// Context contains info about primary object in module
type Context struct {
	Id         uint64
	Level      ContextLevel
	InstanceId uint64
	Path       string
	Depth      int
	Locked     int
}

// CourseContext returns Context record for specific course
func CourseContext(db database.Database, cid uint64) (ctx Context, err error) {
	query := fmt.Sprintf(`SELECT * FROM mdl_context WHERE contextlevel = %d AND
		 instanceid = %d`, CONTEXT_COURSE, cid)

	err = db.QueryRow(query).Scan(&ctx.Id,
		&ctx.Level,
		&ctx.InstanceId,
		&ctx.Path,
		&ctx.Depth,
		&ctx.Locked,
	)
	return
}

// FrozenContexts retruns slice of Contexts which have their "locked" field on.
func FrozenContexts(db database.Database, level ContextLevel) (contexts []Context, err error) {
	query := fmt.Sprintf(`SELECT * FROM mdl_context WHERE contextlevel = %d AND locked > 0`, level)

	var rows *sql.Rows
	rows, err = db.Query(query)
	if err != nil {
		return
	}

	defer rows.Close()

	contexts = make([]Context, 0)
	for rows.Next() {
		var ctx Context
		err = rows.Scan(&ctx.Id,
			&ctx.Level,
			&ctx.InstanceId,
			&ctx.Path,
			&ctx.Depth,
			&ctx.Locked)

		if err != nil {
			return
		}
		contexts = append(contexts, ctx)
	}

	// Check for errors from iterating over rows.
	err = rows.Err()

	return
}

func (ctx Context) String() string {
	return fmt.Sprintf("id: %d \n level: %d, instanceid: %d",
		ctx.Id, ctx.Level, ctx.InstanceId)
}
