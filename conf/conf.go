package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Conf 是配置对象
type Conf struct {
	StorageDriverName string                       `toml:"storage_driver"`
	StorageDrivers    map[string]StorageDriverConf `toml:"storage"`
	ListenAddr        string                       `toml:"listen"`
	Daemon            bool                         `toml:"daemon"`
	LogLevel          string                       `toml:"log_level"`
	LogType           string                       `toml:"log_type"`
	LogFile           string                       `toml:"log_file"`
	PublicURL         string                       `toml:"public_url"`
}

// StorageDriverConf 是每个存储驱动的配置
type StorageDriverConf struct {
	Host     string `toml:"host"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Port     int    `toml:"port"`
	DBName   string `toml:"dbname"`
}

var config Conf

// GetConf 获取配置
func GetConf() Conf {
	return config
}

// ParseConf 解析配置文件
func ParseConf(configFilePath string) Conf {
	if _, err := toml.DecodeFile(configFilePath, &config); err != nil {
		panic(fmt.Sprintf("parse configration file failed: %v", err))
	}

	return config
}
