package orm

import "context"

type Db interface {
	Close()
	Exec(sql string) (lastInsertId int64, rowsAffected int64, err error)
	QueryOne(sql string, result interface{}) error
	Query(sql string, result interface{}) error
	Count(sql string, result interface{}) error
	BeginTx(ctx context.Context) (err error)
	Commit() error
	Rollback() error
	Hook() error
}
