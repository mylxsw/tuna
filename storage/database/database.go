package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mylxsw/tuna/storage"
	log "github.com/sirupsen/logrus"
)

// Register 注册当前驱动到Storage
func Register(driverName, dataSourceName string) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(fmt.Sprintf("%s: can not connect to db: %v", driverName, err))
	}

	storage.Register(driverName, &Storage{
		db:             db,
		driverName:     driverName,
		dataSourceName: dataSourceName,
	})
}

// Storage 使用MySQL为存储引擎
type Storage struct {
	db             *sql.DB
	driverName     string
	dataSourceName string
}

// Set 方法用于在MySQL中创建一个hash与url的对应关系
func (s *Storage) Set(hash, url string, expire int64) (string, error) {

	if expire > 0 {
		// 过期时间设置为“当前时间戳+过期时间”
		expire = time.Now().Unix() + expire
	}

	setSQL := "INSERT INTO tuna_urls (url, expire, hash) VALUES(?, ?, ?)"
	if s.Get(hash) != "" {
		setSQL = "UPDATE tuna_urls SET url = ?, expire = ? WHERE hash = ?"
	}

	stmt, err := s.db.Prepare(setSQL)
	if err != nil {
		return "", fmt.Errorf("Prepare Error: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(url, expire, hash)
	if err != nil {
		return "", fmt.Errorf("Set Failed: %v", err)
	}

	return hash, nil
}

// Get 方法用于在MySQL中查询hash与url的对应关系
func (s *Storage) Get(hash string) string {
	var url string
	err := s.db.QueryRow("SELECT url FROM tuna_urls WHERE hash=? AND (expire >= ? OR expire = 0) LIMIT 1", hash, time.Now().Unix()).Scan(&url)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		log.Warning("%s", err)
	default:
		return url
	}

	return ""
}

// Count 获取当前有多少url
func (s *Storage) Count() int {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM tuna_urls WHERE (expire >= ? or expire = 0)", time.Now().Unix()).Scan(&count)
	if err != nil {
		log.Warning("%s", err)

		return 0
	}
	return count
}

// InitTableForSQLite 用于初始化SQLite表结构
func InitTableForSQLite(dataSourceName string) {
	tableCreateSQL := `
	CREATE TABLE IF NOT EXISTS tuna_urls (
		hash TEXT PRIMARY KEY,
		url TEXT,
		expire INTEGER
	)
	`
	initTable(tableCreateSQL, "sqlite3", dataSourceName)
}

// InitTableForMySQL 用于初始化MySQL表结构
func InitTableForMySQL(dataSourceName string) {
	tableCreateSQL := `
	CREATE TABLE IF NOT EXISTS tuna_urls (
		hash varchar(32) primary key,
		url varchar(255) not null,
		expire int(11) unsigned default 0
	) ENGINE=InnoDB DEFAULT CHARSET=utf8
	`

	initTable(tableCreateSQL, "mysql", dataSourceName)
}

func initTable(tableCreateSQL, name, dataSourceName string) {

	db, err := sql.Open(name, dataSourceName)
	if err != nil {
		log.Fatalf("Create db file failed: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(tableCreateSQL)
	if err != nil {
		log.Fatalf("Create table failed: %v", err)
	}
}
