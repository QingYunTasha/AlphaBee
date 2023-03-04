package usecase

type DispatchMethod string

var (
	Regex DispatchMethod = "regex"
	Glob  DispatchMethod = "glob"
)

type TaskDispatcher struct {
	Method DispatchMethod
}

func (d *TaskDispatcher) Dispatch() {

}
