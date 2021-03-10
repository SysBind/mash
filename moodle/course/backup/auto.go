// moodle/course/backup/auto.go

package backup

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/sysbind/moodle-automated-course-backup/moodle/config"
	"github.com/sysbind/moodle-automated-course-backup/moodle/database"
)

// PreFlight checks some basic settings to see if autobackup should run
func PreFlight(cfg config.Config) (err error) {
	isActive := cfg.GetPluginConf("backup", "backup_auto_active")

	var active int
	if active, err = strconv.Atoi(isActive); err == nil && active == 0 {
		err = errors.New("Automated Course Backups not enabled")
	}

	return
}

// Run starts the automated backup process
func Run(cfg config.Config) (err error) {
	var ids []int

	if ids, err = getCourses(cfg); err != nil {
		log.Fatal(err)
	}

	for i := range ids {
		fmt.Printf("backing up %d \n", ids[i])

		cmd := exec.Command("php", "admin/cli/automated_backup_single.php", strconv.Itoa(ids[i]))
		out, err := cmd.CombinedOutput()
		if err != nil {
			if out != nil {
				fmt.Printf("combined out:\n%s\n", string(out))
			}
			log.Fatalf("cmd.Run() failed with %v\n", err)
		}
	}

	return
}

func getCourses(cfg config.Config) (ids []int, err error) {
	query := fmt.Sprintf("SELECT id FROM mdl_course ORDER BY id DESC")
	var db database.Database = cfg.DB()

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
