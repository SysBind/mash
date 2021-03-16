package pubsub

type Topic int

const (
	AUTO_BACKUP Topic = 0
)

func (t Topic) String() string {
	switch t {
	case AUTO_BACKUP:
		return "Auto-Backup"
	}
	return "Unknown-Topic"
}

type AutoBackupStep int

const (
	CHECK_SKIP AutoBackupStep = 0
	EXECUTE
	CLEANUP
	// final states
	SKIPPED
	COMPLETE
	ERROR
)

type AutoBackupMsg struct {
	Course int
	Step   AutoBackupStep
}
