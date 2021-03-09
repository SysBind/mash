// moodle/config
package config

import (
	"testing"

	"github.com/matryer/is"
)

func TestParseConfig(t *testing.T) {
	is := is.New(t)

	cfg, err := Parse("../testdata/config.php")

	is.NoErr(err)                        // parse config without error
	is.True(cfg.DriverName() == "mysql") // correctly parse db type
}
