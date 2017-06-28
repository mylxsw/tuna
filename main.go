package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mylxsw/tuna/handler"
	mw "github.com/mylxsw/tuna/middleware"
	dbStorage "github.com/mylxsw/tuna/storage/database"
	redisStorage "github.com/mylxsw/tuna/storage/redis"
	redis "gopkg.in/redis.v5"
)

func main() {

	var driverName = "sqlite3"

	switch driverName {
	case "sqlite3":
		sqliteDB := "./test.db"
		// 注册SQLite驱动
		dbStorage.Register("sqlite3", sqliteDB)
		dbStorage.InitTableForSQLite(sqliteDB)
	case "mysql":
		mysqlDataSource := "root:root@tcp(127.0.0.1:3306)/tuna?charset=utf8&parseTime=True&loc=Local"

		dbStorage.Register("mysql", mysqlDataSource)
		dbStorage.InitTableForMySQL(mysqlDataSource)
	case "redis":
		// 注册Redis驱动
		redisStorage.Register("redis", &redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		})
	default:
		panic("no storage driver specified")
	}

	r := mux.NewRouter()

	r.HandleFunc("/", mw.Handler(handler.Welcome, mw.WithHTMLResponse)).Methods("GET")
	r.HandleFunc("/", mw.Handler(handler.Create, mw.WithJSONResponse)).Methods("POST")
	r.HandleFunc("/{hash}", mw.Handler(handler.Query, mw.WithHTMLResponse)).Methods("GET")

	srv := &http.Server{
		Addr:    "127.0.0.1:55555",
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}
