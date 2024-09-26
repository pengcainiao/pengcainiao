package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/martian/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func NewMysql2(host string, db string, maxCon int) *gorm.DB {
	return NewMysql(fmt.Sprintf("%s/%s?charset=utf8mb4&parseTime=true&loc=Local", host, db), maxCon)
}

type logWriter struct {
}

func (*logWriter) Printf(msg string, args ...interface{}) {
	log.Infof(msg, args)
	//qyweixin.ReportErrToWeiXin("gorm_mysql", "mysql error", fmt.Sprintf(msg, args))
}

func newWriter() logger.Writer {
	return &logWriter{}
}

func NewMysql(dsn string, maxCon int) *gorm.DB {
	newLogger := newLogger(
		newWriter(),
		//LogFormatter, //sLog.New(os.Stdout, "\r\n", sLog.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Millisecond * 500, // 慢 SQL 閾值, 0.5 秒
			LogLevel:                  logger.Warn,            // Log level
			IgnoreRecordNotFoundError: true,
			Colorful:                  false, // 禁用彩色列印
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		fmt.Println("gorm.Open error dsn", err.Error())
		return nil
		//panic(fmt.Sprintf("Got errors when connect database, the errors is '%v'", err))
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("gorm.Open error ", err.Error())
		return nil
		//panic(err)
	}
	idle := maxCon
	if maxCon/3 > 10 {
		idle = maxCon / 3
	}
	sqlDB.SetConnMaxIdleTime(time.Minute * 30)
	// 設定空閒連線池中連線的最大數量
	sqlDB.SetMaxIdleConns(idle)
	// 設定開啟資料庫連線的最大數量
	sqlDB.SetMaxOpenConns(maxCon)

	return db
}

// nLogger
type nLogger struct {
	defaultLogger logger.Interface
}

func newLogger(writer logger.Writer, config logger.Config) logger.Interface {
	return &nLogger{
		defaultLogger: logger.New(writer, config),
	}
}

func (l *nLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	return l.defaultLogger.LogMode(logLevel)
}

// Print format & print log
func (l *nLogger) Info(ctx context.Context, msg string, values ...interface{}) {
	l.defaultLogger.Info(ctx, msg, values...)
}

// Print format & print log
func (l *nLogger) Warn(ctx context.Context, msg string, values ...interface{}) {
	l.defaultLogger.Warn(ctx, msg, values...)
	contents := make([]string, 1+len(values))
	contents[0] = msg
	for i, value := range values {
		contents[1+i] = JsonMarshalString(value)
	}
	//qyweixin.ReportErrToWeiXin("gorm_mysql", "mysql warn", contents...)
}

// Print format & print log
func (l *nLogger) Error(ctx context.Context, msg string, values ...interface{}) {
	l.defaultLogger.Error(ctx, msg, values...)
	contents := make([]string, 1+len(values))
	contents[0] = msg
	for i, value := range values {
		contents[1+i] = JsonMarshalString(value)
	}
	//qyweixin.ReportErrToWeiXin("gorm_mysql", "mysql error", contents...)
}

// Print format & print log
func (l *nLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	l.defaultLogger.Trace(ctx, begin, fc, err)
}

func JsonMarshalString(in interface{}) string {
	a, e := json.Marshal(in)
	if e != nil {
		log.Errorf("JsonMarshalString error %s", e.Error())
	}
	return string(a)
}
