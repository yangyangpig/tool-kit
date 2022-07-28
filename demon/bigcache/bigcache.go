package bigcache

import (
	"log"

	"github.com/allegro/bigcache/v3"
)

type BigCacheWarrper struct {
	bigCacheInstance *bigcache.BigCache
}

func initbigcache(cfg bigcache.Config) *BigCacheWarrper {

	bigcache, err := bigcache.NewBigCache(cfg)
	if err != nil {
		log.Fatalf("initbigcache base conf %+v happen error : %+v", cfg, err)
	}

	return &BigCacheWarrper{bigCacheInstance: bigcache}

}

func (s *BigCacheWarrper) Set(key string, data []byte) error {
	return s.bigCacheInstance.Set(key, data)
}

func (s *BigCacheWarrper) Get(key string) ([]byte, error) {
	data, err := s.bigCacheInstance.Get(key)
	if err != nil {
		log.Printf("get data by key %s happen error %+v", key, err)
		return nil, err
	}

	return data, nil
}
