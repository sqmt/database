package database

import (
    "fmt"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func mysqlDriver(o *Option) *gorm.DB {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", o.Username, o.Password, o.Host, o.Port, o.DbName, o.Args)
    mysqlConfig := mysql.Config{
        DSN:                      dsn,  // DSN data source name
        DefaultStringSize:        256,  // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
        DisableDatetimePrecision: true, // disable datetime precision support, which not supported before MySQL 5.6
        // DefaultDatetimePrecision: &datetimePrecision, // default datetime precision
        DontSupportRenameIndex:    true,  // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
        DontSupportRenameColumn:   true,  // use change when rename column, rename rename not supported before MySQL 8, MariaDB
        SkipInitializeWithVersion: false, // smart configure based on used version
    }
    if db, err := gorm.Open(mysql.New(mysqlConfig), GetGormConfig(o)); err != nil {
        return nil
    } else {
        return db
    }
}
