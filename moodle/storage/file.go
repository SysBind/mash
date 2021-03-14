// moodle/storage/file.go

package storage

// StoredFile represents file store under moodledata
type storedFile struct {
	id          int
	filename    string
	contenthash string
	contextid   int
	comonent    string
}

// func StoredFile(
