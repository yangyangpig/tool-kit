package threading

var SinglePlaceholder SinglePlaceholderType
type SinglePlaceholderType = struct {}


type TaskRunner struct {
	limitChan chan SinglePlaceholderType
}

func NewTaskRunner(concurrency int) *TaskRunner  {
	return &TaskRunner{
		limitChan: make(chan SinglePlaceholderType, concurrency),
	}
}

// 如果concurrency=1,可以保证goroutine顺序执行,否则，无法保证
func (tr *TaskRunner) Schedule(task func()) {
	tr.limitChan <- SinglePlaceholder

	go func() {
		defer Recover(func() {
			<-tr.limitChan
		})
	}()

	task()
}

