package database

import (
    "fmt"

    "gorm.io/driver/sqlserver"
    "gorm.io/gorm"
)

func sqlServerDriver(o *Option) *gorm.DB {
    dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", o.Username, o.Password, o.Host, o.Port, o.DbName)
    if db, err := gorm.Open(sqlserver.Open(dsn), GetGormConfig(o)); err != nil {
        return nil
    } else {
        return db
    }
}
