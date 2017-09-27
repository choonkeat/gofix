package gofix

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

type fn func(string, ...interface{}) interface{}

// Use is a row in DB
func Use(t *testing.T, tx *sql.Tx) fn {
	return func(tablename string, args ...interface{}) interface{} {
		var primaryKey string
		var cols, placeholders []string
		var vals []interface{}

		if len(args)%2 != 0 {
			primaryKey = args[0].(string)
			args = args[1:]
		}

		for i, x := range args {
			if i%2 == 0 {
				cols = append(cols, string(x.(string)))
			} else {
				vals = append(vals, x)
				placeholders = append(placeholders, fmt.Sprintf("$%d", len(vals)))
			}
		}

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			tablename,
			strings.Join(cols, ","),
			strings.Join(placeholders, ","),
		)
		if primaryKey != "" {
			query = query + " RETURNING " + primaryKey
			var pkvalue string
			if err := tx.QueryRow(query, vals...).Scan(&pkvalue); err != nil {
				t.Fatalf("failed %#v %#v: %s", query, vals, err.Error())
			}
			return pkvalue
		}

		if _, err := tx.Exec(query, vals...); err != nil {
			t.Fatalf("failed %#v %#v: %s", query, vals, err.Error())
		}
		return ""
	}
}
