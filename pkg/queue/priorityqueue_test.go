package queue

import (
	"container/heap"
	"testing"
)

func Test_PriorityQueue(t *testing.T)  {
	items := map[string]int {
		"one": 3, "tow":2, "three":4,
	}

	pg := make(PriorityQueue, len(items))
	i :=0
	for value, priority := range items {
		pg[i] = &Item{
			value: value,
			priority: priority,
			index: i,
		}
		i++
	}

	heap.Init(&pg)

	// newItem
	newItem := &Item{
		value: "orange",
		priority: 1,
	}
	heap.Push(&pg, newItem)

	// 对指定的item更新优先级
	pg.UpdatePriority(newItem, "four", 8)
	for pg.Len() > 0 {
		item := pg.Pop().(*Item)
		t.Logf("pop value: priority: index %s:%d:%d", item.value, item.priority, item.index)
	}
}