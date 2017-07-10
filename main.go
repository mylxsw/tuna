package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
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
	log "github.com/sirupsen/logrus"
	redis "gopkg.in/redis.v5"
)

var configFilePath string

// Conf 是配置对象
type Conf struct {
	StorageDriverName string                       `toml:"storage_driver"`
	StorageDrivers    map[string]StorageDriverConf `toml:"storage"`
	ListenAddr        string                       `toml:"listen"`
	Daemon            bool                         `toml:"daemon"`
	LogLevel          string                       `toml:"log_level"`
	LogType           string                       `toml:"log_type"`
	LogFile           string                       `toml:"log_file"`
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
	process, err := os.StartProcess(binary, os.Args, &os.ProcAttr{Dir: "", Env: nil, Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}, Sys: nil})
	if err != nil {
		fmt.Println("failed to start process:", err)
		os.Exit(2)
	}

	log.Debugf("start daemon process, pid=%d", process.Pid)
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
			panic(fmt.Sprintf("redis配置错误，dbname必须为数字: %v", err))
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
	defer func() {
		log.Debug("http server stopped")
	}()
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

	log.Debugf("start http server, listen http://%s", config.ListenAddr)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Warningf("server closed because %s", r)
		}
	}()

	flag.Usage = func() {
		fmt.Print("Options:\n\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&configFilePath, "conf", "/etc/tuna.toml", "配置文件路径")
	flag.Parse()

	// 解析配置文件
	var config Conf
	if _, err := toml.DecodeFile(configFilePath, &config); err != nil {
		panic(fmt.Sprintf("parse configration file failed: %v", err))
	}

	// 守护进程模式
	if config.Daemon && os.Getppid() != 1 {
		daemonMode()
	}

	// 配置日志处理
	if config.LogLevel != "" {
		level, err := log.ParseLevel(config.LogLevel)
		if err != nil {
			panic(err)
		}

		log.SetLevel(level)
	}

	if config.LogType == "file" && config.LogFile != "" {
		f, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}

		log.SetOutput(f)
	}

	ctx, cancel := context.WithCancel(context.Background())
	signalHandler(cancel)

	// 初始化驱动配置
	initStorageDriver(config)

	// 注册路由规则
	r := mux.NewRouter()

	type RouteHandler struct {
		Path        string
		HandlerFunc http.HandlerFunc
		Methods     []string
	}

	var routeHandlers = []RouteHandler{
		RouteHandler{
			Path:        "/",
			HandlerFunc: mw.Handler(handler.Welcome, mw.WithHTMLResponse),
			Methods:     []string{"GET"},
		},
		RouteHandler{
			Path:        "/",
			HandlerFunc: mw.Handler(handler.Create, mw.WithJSONResponse),
			Methods:     []string{"POST"},
		},
		RouteHandler{
			Path:        "/{hash}",
			HandlerFunc: mw.Handler(handler.Query, mw.WithHTMLResponse),
			Methods:     []string{"GET"},
		},
	}

	for _, handler := range routeHandlers {
		r.HandleFunc(handler.Path, handler.HandlerFunc).Methods(handler.Methods...)
	}

	// 用于获取所有的API
	r.HandleFunc("/probe/routes", mw.Handler(func(w http.ResponseWriter, r *http.Request) {
		results := make(map[string]map[string]map[string]string)
		results["v1"] = make(map[string]map[string]string)

		for _, handler := range routeHandlers {
			if _, ok := results["v1"][handler.Path]; !ok {
				results["v1"][handler.Path] = make(map[string]string)
			}

			for _, method := range handler.Methods {
				results["v1"][handler.Path][method] = ""
			}
		}

		results["v1"]["/probe/routes"] = map[string]string{"GET": ""}

		jsonRes, _ := json.Marshal(results)
		w.Write(jsonRes)
	}, mw.WithJSONResponse)).Methods("GET")

	// 创建 http server
	startHTTPServer(ctx, config, r)
}
