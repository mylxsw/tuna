package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/mylxsw/tuna/handler"
	mw "github.com/mylxsw/tuna/middleware"
	redisStorage "github.com/mylxsw/tuna/storage/redis"
	redis "gopkg.in/redis.v5"
)

func main() {

	// 注册Redis驱动
	redisStorage.Register(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

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
