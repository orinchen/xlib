package xgorm

import (
	"fmt"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"strings"
	"time"
)

// Config
// 配置项目
type Config struct {
	Name     string `json:"name,optional,default=default"`
	Driver   string `json:"driver,optional,default=mysql" toml:"driver"`
	DBName   string `json:"dbName" toml:"dbName"`
	User     string `json:"user" toml:"user"`
	Password string `json:"password" toml:"password"`
	Host     string `json:"host" toml:"host"`
	Port     int    `json:"port,optional,default=3306" toml:"port"`
	Net      string `json:"net,optional,default=tcp" toml:"net"` // MYSQL 用
	TimeZone string `json:"timezone,optional" toml:"timezone"`   // postgres 用
	OpenSSL  bool   `json:"openssl,optional" toml:"openssl"`     // postgres 用

	// Debug开关
	Debug bool `json:"debug,optional,default=false" toml:"debug"`
	// 最大空闲连接数
	MaxIdleConns int `json:"maxIdleConns,optional,default=10" toml:"maxIdleConns"`
	// 最大活动连接数
	MaxOpenConns int `json:"maxOpenConns,optional,default=100" toml:"maxOpenConns"`
	// 连接的最大存活时间
	ConnMaxLifetime uint `json:"connMaxLifetime,optional,default=300" toml:"connMaxLifetime"`
	// 空闲连接的最大存活时间
	ConnMaxIdleTime uint `json:"connMaxIdleTime,optional,default=300" toml:"connMaxIdleTime"`
	// 开启链路追踪
	EnableTrace bool `json:"enableTrace,optional,default=false" toml:"enableTrace"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Net:             "tcp",
		Driver:          "mysql",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: 300,
		ConnMaxIdleTime: 300,
		EnableTrace:     false,
		TimeZone:        "Asia/Shanghai",
		OpenSSL:         false,
	}
}

// Build ...
func (config *Config) Build() (db *gorm.DB, err error) {
	var dialector gorm.Dialector
	switch strings.ToLower(config.Driver) {
	case "postgres":
		sslmode := ""
		if !config.OpenSSL {
			sslmode = "sslmode=disable"
		}
		dialector = postgres.New(postgres.Config{
			DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=%s %s",
				config.Host, config.User, config.Password, config.DBName, config.Port, config.TimeZone, sslmode),
			PreferSimpleProtocol: false,
			WithoutReturning:     false,
		})
	case "sqlserver":
		dialector = sqlserver.New(sqlserver.Config{
			DSN: fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
				config.User, config.Password, config.Host, config.Port, config.DBName),
			DefaultStringSize: 256,
		})
	case "mysql":
		fallthrough
	default:
		dialector = mysql.New(mysql.Config{
			DSN: (&mysqlDriver.Config{
				User:                 config.User,
				Passwd:               config.Password,
				Net:                  "tcp",
				Addr:                 fmt.Sprintf("%s:%d", config.Host, config.Port),
				DBName:               config.DBName,
				AllowNativePasswords: true,
				InterpolateParams:    true,
				ParseTime:            true,
				Loc:                  time.Local,
				Params: map[string]string{
					"charset": "utf8mb4",
					"timeout": "30s",
				},
			}).FormatDSN(),
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    false, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   false, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		})
	}

	if db, err = gorm.Open(dialector); err != nil {
		return
	}

	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(config.MaxIdleConns)
	sqlDb.SetMaxOpenConns(config.MaxOpenConns)
	sqlDb.SetConnMaxLifetime((time.Duration(config.ConnMaxLifetime)) * time.Second)
	sqlDb.SetConnMaxIdleTime((time.Duration(config.ConnMaxIdleTime)) * time.Second)
	return
}

// MustBuild 返回数据库实例
func (config *Config) MustBuild() *gorm.DB {
	if db, err := config.Build(); err != nil {
		panic(err)
	} else {
		return db
	}
}
