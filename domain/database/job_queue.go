package domain

type DispatchMethod string

var (
	Regex DispatchMethod = "regex"
	Glob  DispatchMethod = "glob"
)

type JobQueue struct {
	Method DispatchMethod
	Jobs   chan Job
}
