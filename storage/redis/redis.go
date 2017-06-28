package redis

import (
	"fmt"
	"time"

	"github.com/mylxsw/tuna/storage"
	redis "gopkg.in/redis.v5"
)

// Register 注册当前驱动到Storage
func Register(driverName string, options *redis.Options) {
	storage.Register(driverName, &Storage{
		client:     redis.NewClient(options),
		options:    options,
		driverName: driverName,
	})
}

// Storage 使用Redis为存储引擎
type Storage struct {
	driverName string
	options    *redis.Options
	client     *redis.Client
}

// Set 方法用于在Redis中创建一个hash与url的对应关系
func (s *Storage) Set(hash, url string, expire int) (string, error) {
	return s.client.Set(getKey(hash), url, time.Duration(expire)*time.Second).Result()
}

// Get 方法用于在Redis中查询hash与url的对应关系
func (s *Storage) Get(hash string) string {
	res, err := s.client.Get(getKey(hash)).Result()
	if err != nil {
		return ""
	}
	return string(res)
}

// Count 获取当前有多少url
func (s *Storage) Count() int {
	// TODO 慎用，效率太低
	return len(s.client.Keys(getKey("*")).Val())
}

func getKey(hash string) string {
	return fmt.Sprintf("tuna:%s", hash)
}
