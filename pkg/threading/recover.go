package threading

import "log"

func Recover(cleanups ... func())  {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		// maybe you need to change the log function
		log.Printf("recover the pain: %+v",p)
	}
}