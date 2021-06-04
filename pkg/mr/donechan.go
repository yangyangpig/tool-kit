package mr

import "sync"


var SinglePlaceholder SinglePlaceholderType
type (
	SinglePlaceholderType = struct {}
)


type DoneChan struct {
	done chan SinglePlaceholderType
	once sync.Once
}

func NewDoneChan() *DoneChan {
	return &DoneChan{
		done: make(chan SinglePlaceholderType),
	}
}

// 避免重复关闭channel
func (dc *DoneChan) Close() {
	dc.once.Do(func() {
		close(dc.done)
	})
}

func (dc *DoneChan) Done() chan SinglePlaceholderType{
	return dc.done
}



