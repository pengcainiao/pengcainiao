package sqlx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pengcainiao/sqlx"
	"gitlab.com/a16624741591/zero/core/env"
	"gitlab.com/a16624741591/zero/core/logx"
	"gitlab.com/a16624741591/zero/core/timex"
)

type SqlxDB struct {
	db *sqlx.DB
}

func MySQL() *SqlxDB {
	db, err := GetSqlConn("mysql", env.DbDSN)
	if err != nil {
		fmt.Println("err:", err)
	}
	return &SqlxDB{db: db}
}

func (s *SqlxDB) MustBegin(ctx context.Context) *sqlx.Tx {
	return s.db.MustBegin()
}

func (s *SqlxDB) MustBeginTx(ctx context.Context, opts *sql.TxOptions) *sqlx.Tx {
	return s.db.MustBeginTx(ctx, opts)
}

func (s *SqlxDB) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	//return logDuration(ctx, query, func() error {
	return s.db.Select(dest, query, args...)
	//})
}

func (s *SqlxDB) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	//return logDuration(ctx, query, func() error {
	return s.db.Get(dest, query, args...)
	//})
}

func (s *SqlxDB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	r, err := s.db.Exec(query, args...)
	return r, err
}

func (s *SqlxDB) MustExec(ctx context.Context, query string, args ...interface{}) sql.Result {
	r := s.db.MustExec(query, args...)
	return r
}
func (s *SqlxDB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	r, err := s.db.Query(query, args...)
	return r, err
}

func (s *SqlxDB) RawDB() *sqlx.DB {
	return s.db
}

func logDuration(ctx context.Context, query string, dbCommand func() error) error {
	var start = timex.Now()
	var err = dbCommand()
	duration := timex.Since(start)
	if duration > slowThreshold.Load() {
		var buf strings.Builder
		buf.WriteString(query)
		logx.WithContext(ctx).WithDuration(duration).WithMessageType("mysql").Slowf("[MYSQL] slowcall on executing: %s", buf.String())
	}
	return err
}
