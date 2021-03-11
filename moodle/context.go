// moodle/context.go

package moodle

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
	level ContextLevel
}
