package task_pool

import "errors"

func Validate(option *Option) error {
	if option.MaxWorkerNum < 0 {
		return errors.New("init worker num param error")
	}
	return nil
}
