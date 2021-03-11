// moodle/moodle_test.go
package moodle

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestConfigPhp(t *testing.T) {
	is := is.New(t)
	_, err := os.Stat("testdata/config.php")
	is.NoErr(err) // moodle/testdata/config.php should exist, see README
	return
}
