package domain

type DispatchMethod string

var (
	Regex DispatchMethod = "regex"
	Glob  DispatchMethod = "glob"
)

type JobDispatcher struct {
	method DispatchMethod
}

func (d *JobDispatcher) Dispatch() {

}
