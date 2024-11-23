package sqlx

import (
	"io"
	"sync"
	"time"

	"github.com/pengcainiao/sqlx"

	"gitlab.com/a16624741591/zero/core/syncx"
)

var (
	MaxIdleConns = 10               // 最大空闲连接数
	MaxOpenConns = 50               // 最大连接数
	MaxLifetime  = 60 * time.Minute // 连接最大存活时间
)

var connManager = syncx.NewResourceManager()

type pingedDB struct {
	*sqlx.DB
	once sync.Once
}

func getCachedSqlConn(driverName, server string) (*pingedDB, error) {
	val, err := connManager.GetResource(server, func() (io.Closer, error) {
		conn, err := newDBConnection(driverName, server)
		if err != nil {
			return nil, err
		}

		return &pingedDB{
			DB: conn,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*pingedDB), nil
}

func GetSqlConn(driverName, server string) (*sqlx.DB, error) {
	pdb, err := getCachedSqlConn(driverName, server)
	if err != nil {
		return nil, err
	}

	pdb.once.Do(func() {
		err = pdb.Ping()
	})
	if err != nil {
		return nil, err
	}

	return pdb.DB, nil
}

func newDBConnection(driverName, datasource string) (*sqlx.DB, error) {
	//conn, err := sql.Open(driverName, datasource)
	conn, err := sqlx.Connect(driverName, datasource)
	if err != nil {
		return nil, err
	}

	// we need to do this until the issue https://github.com/golang/go/issues/9851 get fixed
	// discussed here https://github.com/go-sql-driver/mysql/issues/257
	// if the discussed SetMaxIdleTimeout methods added, we'll change this behavior
	// 8 means we can't have more than 8 goroutines to concurrently access the same database.
	conn.SetMaxIdleConns(MaxIdleConns)
	conn.SetMaxOpenConns(MaxOpenConns)
	conn.SetConnMaxLifetime(MaxLifetime)

	return conn, nil
}
