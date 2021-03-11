// moodle/context.go

package moodle

import (
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
	id         int
	level      ContextLevel
	instanceid int
	path       string
	depth      int
	locked     int
}

func CourseContext(db database.Database, cid int) (ctx Context, err error) {
	query := fmt.Sprintf(`SELECT * FROM mdl_context WHERE contextlevel = %d AND
		 instanceid = %d`, CONTEXT_COURSE, cid)
	fmt.Printf("CourseContext: %s", query)

	err = db.QueryRow(query).Scan(&ctx.id,
		&ctx.level,
		&ctx.instanceid,
		&ctx.path,
		&ctx.depth,
		&ctx.locked,
	)
	return
}

func (ctx Context) String() string {
	return fmt.Sprintf("id: %d \n level: %d, instanceid: %d",
		ctx.id, ctx.level, ctx.instanceid)
}
