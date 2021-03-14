// moodle/storage/file.go

package storage

import (
	"fmt"
	"log"

	"github.com/sysbind/mash/moodle/config"
	"github.com/sysbind/mash/moodle/database"
)

// StoredFile represents file store under moodledata
type StoredFile struct {
	Id          int
	FileName    string
	ContentHash string
	ContextId   int
	Component   string
	FileArea    string
	TimeCreated int
}

func (file StoredFile) String() string {
	return fmt.Sprintf("%s \t %s", file.FileName,
		file.ContentHash)
}

func (file StoredFile) Delete(cfg config.Config) (err error) {
	var canDelete bool
	var db database.Database = ab.cfg.DB()

	if canDelete, err = file.canDelete(db); err != nil || !canDelete {
		return
	}
	log.Println("Deleting File:", file)

	db.DeleteRecord("files", file.Id)
	return
}

func (file StoredFile) canDelete(db database.Database) (bool, error) {
	var count int

	query := fmt.Sprintf("SELECT COUNT(id) FROM mdl_files WHERE contenthash='%s' WHERE id != %d",
		file.ContentHash,
		file.id)

	err := db.QueryRow(query).Scan(&count)

	if err != nil {
		return false, err
	}
	return count == 0, nil
}
