package storage

import (
	"fmt"
	"sort"
	"sync"
)

// Driver 是所有底层存储方式的基础
type Driver interface {
	Set(hash, url string, expire int) (string, error)
	Get(hash string) string
	Count() int
}

var drivers = make(map[string]Driver)
var driversMu sync.RWMutex

// Register 注册一个存储驱动
func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()

	if driver == nil {
		panic("storage: Register driver is nil")
	}

	if _, dup := drivers[name]; dup {
		panic("storage: Register called twice for driver " + name)
	}

	drivers[name] = driver
}

func unregisterAllDrivers() {
	driversMu.Lock()
	defer driversMu.Unlock()

	drivers = make(map[string]Driver)
}

// Drivers 返回所有已经注册了的驱动
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	var list []string
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// Open 返回指定的驱动
func Open(driverName string) (Driver, error) {
	driversMu.RLock()
	driveri, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("storage: unknown driver %q (forgotten import?)", driverName)
	}

	return driveri, nil
}

// Default 返回第一个注册的驱动
func Default() Driver {
	for _, driver := range drivers {
		return driver
	}

	return nil
}
