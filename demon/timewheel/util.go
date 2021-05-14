package timewheel

import (
	"log"
	"time"
)

var Placeholder PlaceholderType

type (
	GenericType     = interface{}
	PlaceholderType = struct{}
)


func Recover(cleanups ...func()) {
	for _, cleanup := range cleanups {
		cleanup()
	}

	if p := recover(); p != nil {
		log.Fatalf("recover the pain %+v", p)
	}
}

func GoSafe(fn func()) {
	go RunSafe(fn)
}

func RunSafe(fn func()) {
	defer Recover()
	fn()
}

// Use the long enough past time as start time, in case timex.Now() - lastTime equals 0.
var initTime = time.Now().AddDate(-1, -1, -1)

func Now() time.Duration {
	return time.Since(initTime)
}

func Since(d time.Duration) time.Duration {
	return time.Since(initTime) - d
}

func Time() time.Time {
	// fmt.Printf("initTime %+v now %+v\n",initTime, initTime.Add(Now()))

	return initTime.Add(Now())
}




