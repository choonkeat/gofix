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
		tuple := gofix.Use(t, tx)

		facebook := tuple("companies", "id",
			"name", "Facebook Inc",
			"address", "1 Hacker Way, Menlo Park, CA 94025, USA",
			"created_at", time.Now())

		fbhelpdesk := tuple("departments", "id",
			"name", "Helpdesk",
			"company_id", facebook,
			"created_at", time.Now())

		t.Run("verify companies", func(t *testing.T) {
			rows, err := tx.Query("SELECT id, name, address FROM companies")
			if err != nil {
				t.Error(err)
				return
			}
			if want, got := true, rows.Next(); want != got {
				t.Errorf("expect at least 1 record, but has none")
			}

			var id, name, address string
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

			var id, name, companyID string
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