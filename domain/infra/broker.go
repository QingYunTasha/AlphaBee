package domain

type Broker interface {
	PushJob(job Job, workerName string)
}
