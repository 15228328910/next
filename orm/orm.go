package orm

import "context"

type Orm interface {
	Table(table string) Orm
	Model(model interface{}) Orm
	Select(field string) Orm
	Join(join string, data ...interface{}) Orm
	Where(format string, data ...interface{}) Orm
	Find(result interface{}) Orm
	First(result interface{}) Orm
	Last(result interface{}) Orm
	Count(result interface{}) Orm
	Group(group string) Orm
	Having(format string, data ...interface{}) Orm
	Offset(start int64) Orm
	Limit(limit int64) Orm
	Order(order string) Orm
	Delete() Orm
	Updates(map[string]interface{}) Orm
	Update(model interface{}) Orm
	UpdateColumn(column string, value interface{}) Orm
	Exec(format string, data ...interface{}) Orm
	LastInsertId() int64
	RowsAffected() int64
	BeginTx(ctx context.Context) Orm
	Commit() Orm
	Rollback() Orm
	Hook(action string, handler func()) Orm
	Error() error
}
