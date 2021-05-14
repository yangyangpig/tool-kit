package timewheel


type TaskRunner struct {
	limitChan chan PlaceholderType
}

func NewTaskRunner(concurrency int) *TaskRunner {
	return &TaskRunner{
		limitChan: make(chan PlaceholderType, concurrency),
	}
}

func (rp *TaskRunner) Schedule(task func()) {
	rp.limitChan <- Placeholder

	go func() {
		defer Recover(func() {
			<-rp.limitChan
		})

		task()
	}()
}

