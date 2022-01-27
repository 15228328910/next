package orm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type NextOrm struct {
	model        interface{}
	table        string
	join         string
	where        string
	update       string
	delete       string
	group        string
	having       string
	offset       int64
	limit        int64
	order        string
	count        string
	field        string
	sql          string
	action       string
	db           Db
	lastInsertId int64
	rowsAffected int64
	err          error
}

func (n *NextOrm) Error() error {
	return n.err
}

func (n *NextOrm) LastInsertId() int64 {
	return n.lastInsertId
}

func (n *NextOrm) RowsAffected() int64 {
	return n.rowsAffected
}

func (n *NextOrm) Join(join string, data ...interface{}) Orm {
	join = fmt.Sprintf(join, data...)
	n.join = join
	return n
}

func (n *NextOrm) Order(order string) Orm {
	n.order = order
	return n
}

func (n *NextOrm) Table(table string) Orm {
	n.table = table
	return n
}

// Model 通过反射获取Table函数
func (n *NextOrm) Model(model interface{}) Orm {
	kind := reflect.TypeOf(model).Kind()
	var values []reflect.Value
	if kind != reflect.Ptr {
		panic("model must be pointer")
	}
	values = reflect.ValueOf(model).MethodByName("Table").Call([]reflect.Value{})
	if len(values) == 0 {
		return n
	}
	if values[0].CanInterface() {
		n.table = values[0].Interface().(string)
	}
	n.model = model
	return n
}

func (n *NextOrm) Select(field string) Orm {
	n.field = field
	return n
}

func (n *NextOrm) Where(format string, data ...interface{}) Orm {
	where := fmt.Sprintf(format, data...)
	n.where = where
	return n
}

func (n *NextOrm) buildQuery() {
	var cause []interface{}
	// field
	if n.field == "" {
		n.field = "*"

	}
	sql := "select %s from %s"
	cause = append(cause, n.field, n.table)

	// join
	if n.join != "" {
		sql += " " + "%s"
		cause = append(cause, n.join)
	}

	// where
	if n.where != "" {
		sql += " " + "where %s"
		cause = append(cause, n.where)
	}

	// group
	if n.group != "" {
		sql += " " + "group by %s"
		cause = append(cause, n.group)
	}

	// having
	if n.having != "" {
		sql += " " + "having %s"
		cause = append(cause, n.having)
	}

	// order
	if n.order != "" {
		sql += " " + "order by %s"
		cause = append(cause, n.order)
	}

	// offset
	if n.offset != 0 {
		sql += " " + "offset %d"
		cause = append(cause, n.offset)
	}

	// limit
	if n.limit != 0 {
		sql += " " + "limit %d"
		cause = append(cause, n.limit)
	}

	// count
	if n.count != "" {
		sql = fmt.Sprintf("select count(%s) from (%s) as d", n.count, sql)
	}

	n.sql = fmt.Sprintf(sql, cause...)
}

func (n *NextOrm) buildUpdate(data interface{}) {
	sql := "update %s set %s"
	n.update = data.(string)

	var cause []interface{}
	cause = append(cause, n.table, n.update)

	// where
	if n.where != "" {
		sql += " " + "where %s"
		cause = append(cause, n.where)
	}

	n.sql = fmt.Sprintf(sql, cause...)
}

func (n *NextOrm) buildDelete() {
	sql := "delete from %s"
	var cause []interface{}
	cause = append(cause, n.table)

	// where
	if n.where != "" {
		sql += " " + "where %s"
		cause = append(cause, n.where)
	}

	n.sql = fmt.Sprintf(sql, cause...)
}

func (n *NextOrm) Find(result interface{}) Orm {
	kind := reflect.TypeOf(result).Kind()
	if kind != reflect.Ptr {
		panic("param must be pointer")
	}
	n.buildQuery()
	n.query(result)
	return n
}

func (n *NextOrm) First(result interface{}) Orm {
	kind := reflect.TypeOf(result).Kind()
	if kind != reflect.Ptr {
		panic("param must be pointer")
	}
	n.limit = 1
	n.queryOne(result)
	return n
}

func (n *NextOrm) Last(result interface{}) Orm {
	kind := reflect.TypeOf(result).Kind()
	if kind != reflect.Ptr {
		panic("param must be pointer")
	}
	n.limit = 1
	var results []interface{}
	n.queryOne(results)
	if len(results) == 0 {
		return n
	}
	value := reflect.ValueOf(result)
	if value.CanSet() {
		value.Set(reflect.ValueOf(results[len(results)-1]))
	}
	return n
}

func (n *NextOrm) Count(result interface{}) Orm {
	n.count = "*"
	n.buildQuery()
	n.db.Count(n.sql, result)
	return n
}

func (n *NextOrm) Group(group string) Orm {
	n.group = group
	return n
}

func (n *NextOrm) Having(format string, data ...interface{}) Orm {
	having := fmt.Sprintf(format, data...)
	n.having = having
	return n
}

func (n *NextOrm) Offset(start int64) Orm {
	n.offset = start
	return n
}

func (n *NextOrm) Limit(limit int64) Orm {
	n.limit = limit
	return n
}

func (n *NextOrm) Delete() Orm {
	n.buildDelete()
	n.lastInsertId, n.rowsAffected, n.err = n.db.Exec(n.sql)
	return n
}

func (n *NextOrm) Updates(m map[string]interface{}) Orm {
	r := reflect.ValueOf(m).MapRange()
	var data []string
	for r.Next() {
		key := r.Key().Interface()
		value := r.Value().Interface()
		data = append(data, fmt.Sprintf("%s='%v'", key, value))
	}
	n.buildUpdate(strings.Join(data, ","))
	n.lastInsertId, n.rowsAffected, n.err = n.db.Exec(n.sql)
	return n
}

func (n *NextOrm) Update(model interface{}) Orm {
	var (
		tm reflect.Type
		vm reflect.Value
	)
	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		tm = reflect.TypeOf(model).Elem()
		vm = reflect.ValueOf(model).Elem()
	} else {
		tm = reflect.TypeOf(model)
		vm = reflect.ValueOf(model)
	}
	len := tm.NumField()
	var data []string
	for i := 0; i < len; i++ {
		key := tm.Field(i).Tag.Get("column")
		// 如果是零值，则不赋值
		if isBlank(vm.Field(i)) {
			continue
		}
		value := vm.Field(i).Interface()
		data = append(data, fmt.Sprintf("%s='%v'", key, value))
	}
	n.buildUpdate(strings.Join(data, ","))
	n.lastInsertId, n.rowsAffected, n.err = n.db.Exec(n.sql)
	return n
}

func (n *NextOrm) UpdateColumn(column string, value interface{}) Orm {
	data := fmt.Sprintf("%s='%v'", column, value)
	n.buildUpdate(data)
	n.lastInsertId, n.rowsAffected, n.err = n.db.Exec(n.sql)
	return n
}

func (n *NextOrm) queryOne(result interface{}) Orm {
	n.buildQuery()
	n.err = n.db.QueryOne(n.sql, result)
	return n
}

func (n *NextOrm) query(result interface{}) Orm {
	n.buildQuery()
	n.err = n.db.Query(n.sql, result)
	return n
}

func (n *NextOrm) Exec(format string, data ...interface{}) Orm {
	sql := fmt.Sprintf(format, data...)
	n.lastInsertId, n.rowsAffected, n.err = n.db.Exec(sql)
	return n
}

func (n *NextOrm) BeginTx(ctx context.Context) Orm {
	n.err = n.db.BeginTx(ctx)
	return n
}

func (n *NextOrm) Commit() Orm {
	n.err = n.db.Commit()
	return n
}

func (n *NextOrm) Rollback() Orm {
	n.err = n.db.Rollback()
	return n
}

func (n *NextOrm) Hook(action string, handler func()) Orm {
	//TODO implement me
	panic("implement me")
}

func NewNextOrm(db Db) Orm {
	return &NextOrm{
		db: db,
	}
}
