package domain

type WorkerStatus string

var (
	Running     WorkerStatus = "running"
	Idling      WorkerStatus = "idling"
	Initiating  WorkerStatus = "initiating"
	Unavailable WorkerStatus = "unavailable"
)

type WorkerQueue struct {
	ID         uint
	WorkerName string
	//WorkerStatus WorkerStatus
	Jobs []Job
}
