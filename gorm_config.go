package database

import (
    "time"

    "gorm.io/gorm"
    "gorm.io/gorm/schema"
)

func GetGormConfig(o *Option) *gorm.Config {
    config := &gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true,
        NamingStrategy: schema.NamingStrategy{
            TablePrefix:   o.Prefix,        // table name prefix, table for `User` would be `t_users`
            SingularTable: o.SingularTable, // use singular table name, table for `User` would be `user` with this option enabled
        },
        PrepareStmt:          true,
        DryRun:               o.DryRun,
        DisableAutomaticPing: !o.AutomaticPing,
    }

    if o.LocalTime {
        config.NowFunc = func() time.Time {
            return time.Now().Local()
        }
    }
    if o.Logger.Open {
        config.Logger = o.logger
    }

    return config
}
