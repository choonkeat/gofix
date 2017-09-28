package gofix_test

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/choonkeat/gofix"
	_ "github.com/lib/pq"
)

func TestUse(t *testing.T) {
	err := withTx(t, func(tx *sql.Tx) error {
		insert := gofix.Use(t, tx)

		facebook := insert("companies", "id",
			"name", "Facebook Inc",
			"address", "1 Hacker Way, Menlo Park, CA 94025, USA",
			"created_at", time.Now())

		fbhelpdesk := insert("departments", "id",
			"name", "Helpdesk",
			"company_id", facebook,
			"created_at", time.Now())

		employee := insert("employees", "id",
			"department_id", fbhelpdesk,
			"name", "Chris",
			"email", "chris@example.com",
			"created_at", time.Now())

		insert("departments_employees",
			"department_id", fbhelpdesk,
			"employee_id", employee)

		t.Run("verify companies", func(t *testing.T) {
			rows, err := tx.Query("SELECT id, name, address FROM companies")
			if err != nil {
				t.Error(err)
				return
			}
			if want, got := true, rows.Next(); want != got {
				t.Errorf("expect at least 1 record, but has none")
			}

			var id int64
			var name, address string
			rows.Scan(&id, &name, &address)
			if want, got := facebook, id; want != got {
				t.Errorf("wanted id=%#v but was %#v", want, got)
			}
			if want, got := "Facebook Inc", name; want != got {
				t.Errorf("wanted name=%#v but was %#v", want, got)
			}
			if want, got := "1 Hacker Way, Menlo Park, CA 94025, USA", address; want != got {
				t.Errorf("wanted address=%#v but was %#v", want, got)
			}
			if want, got := false, rows.Next(); want != got {
				t.Errorf("expect only 1 record, but has more")
			}
		})

		t.Run("verify departments", func(t *testing.T) {
			rows, err := tx.Query("SELECT id, name, company_id FROM departments")
			if err != nil {
				t.Error(err)
				return
			}
			if want, got := true, rows.Next(); want != got {
				t.Errorf("expect at least 1 record, but has none")
			}

			var id, companyID int64
			var name string
			rows.Scan(&id, &name, &companyID)
			if want, got := fbhelpdesk, id; want != got {
				t.Errorf("wanted id=%#v but was %#v", want, got)
			}
			if want, got := "Helpdesk", name; want != got {
				t.Errorf("wanted name=%#v but was %#v", want, got)
			}
			if want, got := facebook, companyID; want != got {
				t.Errorf("wanted company_id=%#v but was %#v", want, got)
			}
			if want, got := false, rows.Next(); want != got {
				t.Errorf("expect only 1 record, but has more")
			}
		})

		t.Run("verify departments_employees", func(t *testing.T) {
			rows, err := tx.Query("SELECT department_id, employee_id FROM departments_employees")
			if err != nil {
				t.Error(err)
				return
			}
			if want, got := true, rows.Next(); want != got {
				t.Errorf("expect at least 1 record, but has none")
			}

			var department_id, employee_id int64
			rows.Scan(&department_id, &employee_id)
			if want, got := fbhelpdesk, department_id; want != got {
				t.Errorf("wanted department_id=%#v but was %#v", want, got)
			}
			if want, got := employee, employee_id; want != got {
				t.Errorf("wanted employee_id=%#v but was %#v", want, got)
			}
		})

		return nil
	})
	if err != nil {
		t.Fatal(err.Error())
	}
}

func withTx(t *testing.T, callback func(*sql.Tx) error) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	return callback(tx)
}
