package container

import (
	"container/list"
	"crypto/md5"
	"encoding/hex"
)

type (
	// hap map 的value的list用于解决hash冲突
	HashMap struct {
		Element map[interface{}]*list.List
	}
)

func NewHapMap() *HashMap {
	return &HashMap{}
}

func (h *HashMap)Add(key string, value interface{})  {
	// hash函数
	hk := md5.Sum([]byte(key))

	hashKey := hex.EncodeToString(hk[:])
	if h.Element[hashKey] != nil {
		// 冲突了
		h.Element[hashKey].PushBack(value)
		return
	}
	l := list.New()
	l.PushFront(value)
	h.Element[hashKey] = l
	return
}

func (h *HashMap) Delete(key string) {

}
