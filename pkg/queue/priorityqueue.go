package queue

import "container/heap"

/*
* 优先队列：在队列中，可以任意顺序添加，在取出数据时，首先选择最小值
* 堆规则：子类数字总是大于其父类数字
* 数字添加到末尾，如果父类数字较大，子类与父类交互
* 数字取出，从堆顶取出数字，在堆中，最小值保存在顶部位置，取出栈顶最小值，需要重新组织堆的
* 结构，从结尾的数字移动到顶部，子类数字比其父类小，此时，相邻子类中数字较小的与父类数字交互
* 重复此操作，直到不发生替换。
*/
type (
	Item struct {
		value string
		priority int
		index int
	}

	PriorityQueue []*Item
)

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Swap(i, j int)  {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	item := old[len(*pq)-1]

	*pq = old[0:len(*pq)-1]

	return item
}

func (pq *PriorityQueue) UpdatePriority(item *Item, value string, priority int) {
	item.priority = priority
	item.value = value
	heap.Fix(pq, item.index)
}
