package mr

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
)

var errDummy = errors.New("dummy")

var mockTestData = []struct{
	mapper      MapperFunc
	reducer     ReducerFunc
	expectErr   error
	expectValue interface{}
} {
	{
	//	mapper: func(item interface{}, write Write, cancel func(error)) {
	//	v := item.(int)
	//	if v % 3 == 0 { // 3正除
	//		cancel(errDummy)
	//	}
	//	write.Pipe(v * v)
	//},
		mapper: nil,
		reducer: nil,
		expectErr: nil,
		expectValue: 30,
	},

	{
		//mapper: func(item interface{}, write Write, cancel func(error)) {
		//	v := item.(int)
		//	if v % 3 == 0 { // 3正除
		//		//cancel(errDummy)
		//	}
		//	write.Pipe(v * v)
		//},
		mapper: nil,
		//reducer: func(pipe <-chan interface{}, write Write, cancel func(error)) {
		//	var result int
		//	for item := range pipe {
		//		result += item.(int)
		//		fmt.Printf("reducer here %+v\n", result)
		//		//if result > 10 {
		//		//	cancel(errDummy)
		//		//}
		//	}
		//	write.Pipe(result)
		//},
		reducer: nil,
		expectErr: nil,
		expectValue: 30,
	},

}

func TestMapReduce(t *testing.T) {
	var value uint32
	for k, v := range mockTestData {
		tag := fmt.Sprintf("test_index_%d", k)
		t.Run(tag, func(t *testing.T) {
			if v.mapper == nil {
				v.mapper = func(item interface{}, write Write, cancel func(error)) {
					v := item.(int)
					write.Pipe(v * v)
				}
			}
			if v.reducer == nil {
				v.reducer = func(pipe <-chan interface{}, write Write, cancel func(error)) {
					for item := range pipe {
						atomic.AddUint32(&value, uint32(item.(int)))
					}
					write.Pipe(atomic.LoadUint32(&value))
				}
			}

			if v.mapper != nil && v.reducer != nil {
				resp, err := MapReduce(func(source chan<- interface{}) {
					for i:=1; i < 5; i++ {
						source <- i
					}
				}, v.mapper, v.reducer, WithRoutines(2))
				if err != nil {
					t.Errorf("MapReduce happen error %+v", err)
					return
				}
				//time.Sleep(time.Second * 3)
				t.Logf("MapReduce target value: %d real value %d", 5, resp)
			}
		})
	}
}
