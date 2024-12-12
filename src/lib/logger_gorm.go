package lib

import (
	"context"
	"errors"
	"fmt"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

/*
GormLogger logger for gorm logging [subbed from main logger]
*/

type GormLogger struct {
	*Logger
	gormlogger.Config
	traceStr, traceErrStr, traceWarnStr string
}

func (l *GormLogger) setup() {
	var (
		traceStr     = "\n```sql [%.6f ms] [rows:%v]\n%s\n```"
		traceWarnStr = "%s\n```sql [%.6f ms] [rows:%v]\n%s\n```"
		traceErrStr  = "%s\n```sql [%.6f ms] [rows:%v]\n%s\n```"
	)

	if l.Config.Colorful {
		traceStr = fmt.Sprintf(
			"\n```sql %s %s\n%s\n%s```",
			gormlogger.Yellow+"[%.6f ms]", // latency
			gormlogger.Blue+"[rows:%v]",   // rows
			gormlogger.Green+"%s",         // query
			gormlogger.Reset,
		)
		traceWarnStr = fmt.Sprintf(
			"%s\n```sql %s %s\n%s\n%s```",
			gormlogger.Yellow+"%s"+gormlogger.Reset,
			gormlogger.Red+"[%.6f ms]",     // latency
			gormlogger.Magenta+"[rows:%v]", // rows
			gormlogger.Blue+"%s",           // query
			gormlogger.Reset,
		)
		traceErrStr = fmt.Sprintf(
			"%s\n```sql %s %s\n%s\n%s```",
			gormlogger.Red+"%s"+gormlogger.Reset,
			gormlogger.Yellow+"[%.6f ms]", // latency
			gormlogger.Green+"[rows:%v]",  // rows
			gormlogger.Magenta+"%s",       // query
			gormlogger.Reset,
		)
	}

	l.traceStr = traceStr
	l.traceErrStr = traceErrStr
	l.traceWarnStr = traceWarnStr
}

/* --- GORM logger interface implementation --- */

/*
LogMode set log mode
*/

func (l GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	l.LogLevel = level
	return &l
}

/*
Info prints info messages
*/

func (l GormLogger) Info(_ context.Context, str string, args ...any) {
	if l.LogLevel >= gormlogger.Info {
		l.Debugf(str, args...)
	}
}

/*
Warn prints warn messages
*/

func (l GormLogger) Warn(_ context.Context, str string, args ...any) {
	if l.LogLevel >= gormlogger.Warn {
		l.Warnf(str, args...)
	}
}

/*
Error prints error messages
*/

func (l GormLogger) Error(_ context.Context, str string, args ...any) {
	if l.LogLevel >= gormlogger.Error {
		l.Errorf(str, args...)
	}
}

/*
Trace prints trace messages
*/

func (l GormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}

	elapsed := time.Since(begin)
	ms := float64(elapsed.Nanoseconds()) / 1e6

	if err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gormlogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError) {
		sql, rows := fc()
		l.SugaredLogger.Errorf(l.traceErrStr, err, ms, rows, sql)
		return
	} else if elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn {
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		l.SugaredLogger.Warnf(l.traceWarnStr, slowLog, ms, rows, sql)
		return
	} else if l.LogLevel >= gormlogger.Info {
		sql, rows := fc()
		l.Debugf(l.traceStr, ms, rows, sql)
		return
	}

}
