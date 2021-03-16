package backup

type CourseBackup struct {
	CourseId  int
	StartTime int
	EndTime   int
}

func GetBackup(cid int) CourseBackup {
	return CourseBackup{}
}
