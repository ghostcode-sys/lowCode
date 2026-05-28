package database

import (
	"context"
	"lowleveldesign/cricbuzz/utility"
	"reflect"
	"time"
)

type T interface{}

func SelectQuery[T any](db *DBConnection, query string, args []any, timeout time.Duration) (res []T, err error) {
	defer utility.HandlePanic()
	res = make([]T, 0)

	if db == nil || db.err != nil {
		err = db.err
		return
	}

	// Executing select query
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := db.conn.QueryContext(ctx, query, args...)

	if err != nil {
		return
	}

	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return
	}
	for rows.Next() {
		var item T
		v := reflect.ValueOf(item).Elem()
		t := v.Type()

		fieldMap := make(map[string]reflect.Value)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tagName := field.Tag.Get("db")

			if tagName != "" {
				fieldMap[tagName] = v.Field(i)
			}

			fieldMap[field.Name] = v.Field(i)
		}

		dest := make([]any, len(columns))

		for i, colName := range columns {
			if field, ok := fieldMap[colName]; ok {
				dest[i] = field.Addr().Interface()
			} else {
				var dummy any
				dest[i] = &dummy
			}
		}

		err = rows.Scan(dest...)
		if err != nil {
			return
		}

		res = append(res, item)
	}

	return
}

func ExecuteQuery(db *DBConnection, query string, args []any, timeout time.Duration) (rowAffected int64, err error) {
	defer utility.HandlePanic()
	rowAffected = 0
	if db == nil || db.err != nil {
		err = db.err
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := db.conn.ExecContext(ctx, query, args...)

	if err != nil {
		return
	}

	rowAffected, err = result.RowsAffected()
	return
}

func CreateTable(db *DBConnection, query string, timeout time.Duration) (err error) {
	defer utility.HandlePanic()

	if db == nil || db.err != nil {
		err = db.err
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = db.conn.ExecContext(ctx, query)

	if err != nil {
		return
	}
	return
}
