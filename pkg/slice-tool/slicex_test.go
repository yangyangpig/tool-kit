package slice_tool

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestElem_RetrieveExpiry(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	nm := rand.Int63n(100)
	contain := NewElem(int(nm))
	for i := 0; i < int(nm); i++ {
		worker := &Worker{
			Fn: func() {
				fmt.Println()
			},
			RecycleTime: time.Now(),
		}
		contain.Insert(worker)
	}

	time.Sleep(time.Second * 10)
	refresh := &Worker{
		Fn: func() {
			fmt.Println("this is refresh worker")
		},
		RecycleTime: time.Now(),
	}
	contain.Insert(refresh)

	expireWork := contain.RetrieveExpiry(time.Second* 10)

	t.Logf("expiry worker %+v", expireWork)
}
