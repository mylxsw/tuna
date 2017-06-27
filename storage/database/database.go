package mysql

import (
	"database/sql"
	"fmt"

	"github.com/mylxsw/tuna/storage"
)

// Register 注册当前驱动到Storage
func Register(dataSourceName string) {
	db, err := sql.Open("database", dataSourceName)
	if err != nil {
		panic(fmt.Sprintf("database: can not connect to mysql: %v", err))
	}

	storage.Register("database", &Storage{
		db: db,
	})
}

// Storage 使用MySQL为存储引擎
type Storage struct {
	db *sql.DB
}

// Set 方法用于在MySQL中创建一个hash与url的对应关系
func (s *Storage) Set(hash, url string, expire int) (string, error) {
	return "", nil
}

// Get 方法用于在MySQL中查询hash与url的对应关系
func (s *Storage) Get(hash string) string {
	stmt, _ := s.db.Prepare("SELECT url FROM tuna_urls WHERE hash=? LIMIT 1")
	rows, _ := stmt.Query(hash)
	rows.Next()
	return ""
}

// Count 获取当前有多少url
func (s *Storage) Count() int {
	return 0
}
