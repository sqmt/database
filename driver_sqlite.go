package database

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func sqliteDriver(o *Option) *gorm.DB {
    if db, err := gorm.Open(sqlite.Open(o.DbName), GetGormConfig(o)); err != nil {
        return nil
    } else {
        return db
    }
}
