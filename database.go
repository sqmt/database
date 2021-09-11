package database

import (
    "time"

    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var drivers map[string]Driver

type Driver func(o *Option) *gorm.DB

type Option struct {
    Driver        string        `json:"driver" yaml:"driver"`
    Host          string        `json:"host" yaml:"host"`
    Port          string        `json:"port" yaml:"port"`
    DbName        string        `json:"dbName" yaml:"dbName"`
    Username      string        `json:"username" yaml:"username"`
    Password      string        `json:"password" yaml:"password"`
    Args          string        `json:"args" yaml:"args"`
    Prefix        string        `json:"prefix" yaml:"prefix"`
    LocalTime     bool          `json:"localTime" yaml:"localTime"`
    DryRun        bool          `json:"dryRun" yaml:"dryRun"`
    AutomaticPing bool          `json:"automaticPing" yaml:"automaticPing"`
    MaxIdleConns  int           `json:"maxIdleConns" yaml:"maxIdleConns"`
    MaxOpenConns  int           `json:"maxOpenConns" yaml:"maxOpenConns"`
    MaxLifetime   time.Duration `json:"maxLifetime" yaml:"maxLifetime"`
    SingularTable bool          `json:"singularTable" yaml:"singularTable"`
    Logger        *Logger       `mapstructure:"logger" json:"logger" yaml:"logger"`
    logger        logger.Interface
}

func (o *Option) NewLogger(writer LoggerWriterAdapter) {
    o.logger = newLogger(o, writer)
}

func init() {
    SetAdapter("mysql", mysqlDriver)
    SetAdapter("sqlite", sqliteDriver)
    SetAdapter("postgres", postgresDriver)
    SetAdapter("sqlserver", sqliteDriver)
}

func SetAdapter(name string, driver Driver) {
    if drivers == nil {
        drivers = map[string]Driver{}
    }
    drivers[name] = driver
}

func New(o *Option) *gorm.DB {
    var db *gorm.DB
    if o.DbName == "" {
        return nil
    }

    if adapter, ok := drivers[o.Driver]; ok {
        db = adapter(o)
    } else {
        panic("database driver " + o.Driver + " not support")
    }
    if db == nil {
        panic("database open failed")
    }
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(o.MaxIdleConns)
    sqlDB.SetMaxOpenConns(o.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(o.MaxLifetime)

    return db
}
