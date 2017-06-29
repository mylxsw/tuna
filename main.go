package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"github.com/mylxsw/tuna/handler"
	mw "github.com/mylxsw/tuna/middleware"
	dbStorage "github.com/mylxsw/tuna/storage/database"
	redisStorage "github.com/mylxsw/tuna/storage/redis"
	redis "gopkg.in/redis.v5"
)

var configFilePath string

// Conf 是配置对象
type Conf struct {
	StorageDriverName string                       `toml:"storage_driver"`
	StorageDrivers    map[string]StorageDriverConf `toml:"storage"`
	ListenAddr        string                       `toml:"listen"`
	Daemon            bool                         `toml:"daemon"`
}

// StorageDriverConf 是每个存储驱动的配置
type StorageDriverConf struct {
	Host     string `toml:"host"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Port     int    `toml:"port"`
	DBName   string `toml:"dbname"`
}

func daemonMode() {
	binary, err := exec.LookPath(os.Args[0])
	if err != nil {
		fmt.Println("failed to lookup binary:", err)
		os.Exit(2)
	}
	_, err = os.StartProcess(binary, os.Args, &os.ProcAttr{Dir: "", Env: nil, Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}, Sys: nil})
	if err != nil {
		fmt.Println("failed to start process:", err)
		os.Exit(2)
	}

	os.Exit(0)
}

func initStorageDriver(config Conf) {
	switch config.StorageDriverName {
	case "sqlite3":
		sqliteDB := config.StorageDrivers["sqlite"].DBName
		// 注册SQLite驱动
		dbStorage.Register("sqlite3", sqliteDB)
		dbStorage.InitTableForSQLite(sqliteDB)
	case "mysql":
		mysqlConf := config.StorageDrivers["mysql"]
		mysqlDataSource := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			mysqlConf.Username,
			mysqlConf.Password,
			mysqlConf.Host,
			mysqlConf.Port,
			mysqlConf.DBName)

		dbStorage.Register("mysql", mysqlDataSource)
		dbStorage.InitTableForMySQL(mysqlDataSource)
	case "redis":
		redisConf := config.StorageDrivers["redis"]
		redisDB, err := strconv.Atoi(redisConf.DBName)
		if err != nil {
			log.Fatalf("redis配置错误，dbname必须为数字: %v", err)
		}
		// 注册Redis驱动
		redisStorage.Register("redis", &redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
			Password: redisConf.Password,
			DB:       redisDB,
		})
	default:
		panic("no storage driver specified")
	}
}

func startHTTPServer(ctx context.Context, config Conf, handler http.Handler) {
	srv := &http.Server{
		Addr:    config.ListenAddr,
		Handler: handler,
	}

	// server关闭
	go func() {
		select {
		case <-ctx.Done():
			srv.Shutdown(ctx)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}

func main() {

	flag.Usage = func() {
		fmt.Print("Options:\n\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&configFilePath, "conf", "/etc/tuna.toml", "配置文件路径")
	flag.Parse()

	// 解析配置文件
	var config Conf
	if _, err := toml.DecodeFile(configFilePath, &config); err != nil {
		log.Fatalf("parse configration file failed: %v", err)
	}

	// 守护进程模式
	if config.Daemon && os.Getppid() != 1 {
		daemonMode()
	}

	ctx, cancel := context.WithCancel(context.Background())
	signalHandler(cancel)

	// 初始化驱动配置
	initStorageDriver(config)

	// 注册路由规则
	r := mux.NewRouter()

	r.HandleFunc("/", mw.Handler(handler.Welcome, mw.WithHTMLResponse)).Methods("GET")
	r.HandleFunc("/", mw.Handler(handler.Create, mw.WithJSONResponse)).Methods("POST")
	r.HandleFunc("/{hash}", mw.Handler(handler.Query, mw.WithHTMLResponse)).Methods("GET")

	// 创建 http server
	startHTTPServer(ctx, config, r)
}
