/*
 *  Author: Orin Chen
 *   Email: orinchen@gmail.com
 *    Time: 2020/1/8 10:35 上午
 * Project: market
 *    File: instances.go
 *     IDE: GoLand
 */

package xgorm

import (
	"gorm.io/gorm"
	"sync"
)

var mutex sync.Mutex
var instances sync.Map
var once sync.Once

const DefaultConnName = "default"

func Db() *gorm.DB {
	return GetDb(DefaultConnName)
}

func Init(configs ...Config) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, config := range configs {
		instances.Store(config.Name, config.Build())
	}
}

func GetDb(name string) *gorm.DB {
	v, hasDb := instances.Load(name)
	if hasDb {
		if db, ok := v.(*gorm.DB); ok {
			return db
		}
	}
	return nil
}
