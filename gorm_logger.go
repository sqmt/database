package database

import (
    "context"
    "errors"
    "fmt"
    "time"

    "gorm.io/gorm/logger"
    "gorm.io/gorm/utils"
)

type Logger struct {
    writer                              LoggerWriterAdapter
    logLevel                            logger.LogLevel
    Open                                bool          `json:"open" yaml:"open"`
    Level                               string        `json:"level" yaml:"level"`
    SlowThreshold                       time.Duration `json:"slowThreshold" yaml:"slowThreshold"`
    IgnoreRecordNotFoundError           bool          `json:"ignoreRecordNotFoundError" yaml:"ignoreRecordNotFoundError"`
    InfoStr, WarnStr, ErrStr            string
    TraceStr, TraceErrStr, TraceWarnStr string
}

type LoggerWriterAdapter interface {
    Info(args ...interface{})
    Warn(args ...interface{})
    Error(args ...interface{})
}

func newLogger(o *Option, writer LoggerWriterAdapter) *Logger {
    if !o.Logger.Open || writer == nil {
        return nil
    }
    l := &Logger{
        writer:                    writer,
        Open:                      o.Logger.Open,
        Level:                     o.Logger.Level,
        IgnoreRecordNotFoundError: o.Logger.IgnoreRecordNotFoundError,
    }
    if o.Logger.SlowThreshold <= 0 {
        l.SlowThreshold = 200 * time.Millisecond
    }
    switch o.Logger.Level {
    case "info":
        l.logLevel = logger.Info
    case "warn":
        l.logLevel = logger.Warn
    case "error":
        l.logLevel = logger.Error
    default:
        l.logLevel = logger.Error
    }

    if o.Logger.InfoStr == "" {
        l.InfoStr = "%s\n[info] "
    }
    if o.Logger.WarnStr == "" {
        l.InfoStr = "%s\n[warn] "
    }
    if o.Logger.ErrStr == "" {
        l.InfoStr = "%s\n[error] "
    }
    if o.Logger.TraceStr == "" {
        l.InfoStr = "%s\n[%.3fms] [rows:%v] %s"
    }
    if o.Logger.TraceWarnStr == "" {
        l.InfoStr = "%s %s\n[%.3fms] [rows:%v] %s"
    }
    if o.Logger.TraceErrStr == "" {
        l.InfoStr = "%s %s\n[%.3fms] [rows:%v] %s"
    }

    return l
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
    c := *l
    c.logLevel = level
    return &c
}

func (l *Logger) Info(ctx context.Context, s string, i ...interface{}) {
    if l.logLevel >= logger.Info {
        l.writer.Info(fmt.Sprintf(l.InfoStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...))
    }
}

func (l *Logger) Warn(ctx context.Context, s string, i ...interface{}) {
    if l.logLevel >= logger.Warn {
        l.writer.Warn(fmt.Sprintf(l.WarnStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...))
    }
}

func (l *Logger) Error(ctx context.Context, s string, i ...interface{}) {
    if l.logLevel >= logger.Error {
        l.writer.Error(fmt.Sprintf(l.WarnStr+s, append([]interface{}{utils.FileWithLineNum()}, i...)...))
    }
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
    if l.logLevel <= logger.Silent {
        return
    }

    elapsed := time.Since(begin)
    switch {
    case err != nil && l.logLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
        sql, rows := fc()
        if rows == -1 {
            l.writer.Error(fmt.Sprintf(l.TraceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql))
        } else {
            l.writer.Error(fmt.Sprintf(l.TraceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql))
        }
    case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.logLevel >= logger.Warn:
        sql, rows := fc()
        slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
        if rows == -1 {
            l.writer.Warn(fmt.Sprintf(l.TraceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql))
        } else {
            l.writer.Warn(fmt.Sprintf(l.TraceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql))
        }
    case l.logLevel == logger.Info:
        sql, rows := fc()
        if rows == -1 {
            l.writer.Info(fmt.Sprintf(l.TraceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql))
        } else {
            l.writer.Info(fmt.Sprintf(l.TraceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql))
        }
    }
}
