package bigcache

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
)

var initBigcache *BigCacheWarrper
var bigcacheConfig = bigcache.Config{}

type mockData struct {
	name string
	age  uint16
}

func init() {
	initBigcache = initbigcache(bigcache.Config{
		Shards:               1024,             // 数据存储空间大小，必须是2的倍数
		LifeWindow:           10 * time.Minute, // 窗口生命周期
		CleanWindow:          5 * time.Minute,  // 清除过期数据条目的窗口间隔
		MaxEntriesInWindow:   1000 * 10 * 60,
		MaxEntrySize:         500,
		StatsEnabled:         false,
		Verbose:              true,
		Hasher:               nil,
		HardMaxCacheSize:     8192, // 最大缓存上限
		OnRemove:             nil,
		OnRemoveWithMetadata: nil,
		OnRemoveWithReason:   nil,
		Logger:               nil,
	})
}

func Test_FirstDemon(t *testing.T) {
	data := mockData{name: "jack", age: 12}
	marshalValue, _ := json.Marshal(data)

	err := initBigcache.Set("first-demon", marshalValue)
	if err != nil {
		log.Fatalf("get data from cache error %+v", err)
	}
	rawdatabyte, err := initBigcache.Get("first-demon")
	if err != nil {
		log.Fatalf("get data from cache error %+v", err)
	}
	log.Printf("get from cache value %+v", rawdatabyte)
	rawdata := &mockData{}
	err = json.Unmarshal(rawdatabyte, rawdata)
	if err != nil {
		log.Fatalf("get data from cache error %+v", err)
	}
	log.Printf("data value %+v", rawdata)
}