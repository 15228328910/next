package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"time"
)

type Mysql struct {
	db *sql.DB
	tx *sql.Tx
}

func (m *Mysql) Hook() error {
	return nil
}

func (m *Mysql) BeginTx(ctx context.Context) (err error) {
	m.tx, err = m.db.BeginTx(ctx, nil)
	return
}

func (m *Mysql) Commit() error {
	return m.tx.Commit()
}

func (m *Mysql) Rollback() error {
	return m.tx.Rollback()
}

func (m *Mysql) Count(sqlStr string, result interface{}) error {
	if reflect.TypeOf(result).Kind() != reflect.Ptr {
		return errors.New("must be pointer")
	}
	elem := reflect.TypeOf(result).Elem()
	data := reflect.New(elem)

	var rows *sql.Row
	if m.tx == nil {
		rows = m.db.QueryRow(sqlStr)
	} else {
		rows = m.tx.QueryRow(sqlStr)
	}
	rows.Scan(data.Interface())
	reflect.Indirect(reflect.ValueOf(result)).Set(data.Elem())
	return rows.Err()
}

func (m *Mysql) QueryOne(sqlStr string, result interface{}) (err error) {
	score := reflect.Indirect(reflect.ValueOf(result))
	var rows *sql.Rows
	if m.tx == nil {
		rows, err = m.db.Query(sqlStr)
	} else {
		rows, err = m.tx.Query(sqlStr)
	}
	if err != nil {
		return
	}
	defer rows.Close()

	// 获取返回值结构体
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	elem := reflect.TypeOf(result).Elem()
	tagValue := make(map[string]reflect.Value, 0)
	item := reflect.New(elem)
	for i := 0; i < elem.NumField(); i++ {
		tagValue[elem.Field(i).Tag.Get("column")] = item.Elem().Field(i)
	}
	fmt.Println("name is:", reflect.TypeOf(result), reflect.TypeOf(result).Elem(), elem.String())

	for rows.Next() {
		var data []interface{}
		for _, name := range columns {
			if v, ok := tagValue[name]; ok {
				data = append(data, v.Addr().Interface())
			} else {
				var empty interface{}
				data = append(data, &empty)
			}
		}
		rows.Scan(data...)
		score.Set(item.Elem())
		break
	}
	return nil
}

func (m *Mysql) Query(sqlStr string, result interface{}) (err error) {
	var rows *sql.Rows
	if m.tx == nil {
		rows, err = m.db.Query(sqlStr)
	} else {
		rows, err = m.tx.Query(sqlStr)
	}
	defer rows.Close()

	// 获取返回值结构体
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	if reflect.Slice != reflect.TypeOf(result).Elem().Kind() {
		return errors.New("result必须为数组")
	}
	tagValue := make(map[string]reflect.Value, 0)

	// 获取数组元素类型
	elem := reflect.TypeOf(result).Elem().Elem()

	score := reflect.Indirect(reflect.ValueOf(result))
	for rows.Next() {
		var (
			data  []interface{}
			vItem reflect.Value
			tItem reflect.Type
		)
		if reflect.Ptr == reflect.TypeOf(result).Elem().Elem().Kind() {
			vItem = reflect.New(elem.Elem())
			tItem = reflect.TypeOf(result).Elem().Elem().Elem()
		} else {
			vItem = reflect.New(elem)
			tItem = reflect.TypeOf(result).Elem().Elem()
		}
		for i := 0; i < tItem.NumField(); i++ {
			tagValue[tItem.Field(i).Tag.Get("column")] = vItem.Elem().Field(i)
		}
		for _, name := range columns {
			if v, ok := tagValue[name]; ok {
				data = append(data, v.Addr().Interface())
			} else {
				var empty interface{}
				data = append(data, &empty)
			}
		}
		rows.Scan(data...)
		score.Set(reflect.Append(score, vItem))
	}
	return nil
}

func (m *Mysql) Exec(sqlStr string) (lastInsertId int64, rowsAffected int64, err error) {
	var result sql.Result
	if m.tx == nil {
		result, err = m.db.Exec(sqlStr)
	} else {
		result, err = m.tx.Exec(sqlStr)
	}
	if err != nil || result == nil {
		return
	}
	lastInsertId, _ = result.LastInsertId()
	rowsAffected, _ = result.RowsAffected()
	return
}

func (m *Mysql) connect(dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	m.db = db
	return m.db.Ping()
}

func (m *Mysql) Close() {
	m.db.Close()
}

func NewMysql(dsn string) (db Db, err error) {
	m := &Mysql{}
	err = m.connect(dsn)
	return m, err
}
