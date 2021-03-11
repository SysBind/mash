// moodle/context_test.go

package moodle

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/sysbind/mash/moodle/config"
	"github.com/sysbind/mash/moodle/database"
)

func TestCourseContext(t *testing.T) {
	is := is.New(t)

	cfg, err := config.Parse("testdata/config.php")
	is.NoErr(err) // parse config without error

	db, err := database.Open(context.Background(), cfg.DriverName(), cfg.DSN())
	is.NoErr(err) // open database
	cfg.SetDatabase(db)

	ctx, err := CourseContext(db, 1)
	is.NoErr(err)                // get course context without error
	is.True(ctx.instanceid == 1) // get context with instancid we asked for
}
