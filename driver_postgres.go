package database

import (
    "fmt"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func postgresDriver(o *Option) *gorm.DB {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s %s", o.Host, o.Username, o.Password, o.DbName, o.Port, o.Args)
    postgresConfig := postgres.Config{
        DSN: dsn, // DSN data source name
    }
    if db, err := gorm.Open(postgres.New(postgresConfig), GetGormConfig(o)); err != nil {
        return nil
    } else {
        return db
    }
}
