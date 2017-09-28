# Experiment to have easier fixtures in Go

``` go
insert := gofix.Use(t, tx)

facebook := insert("companies", "id",
  "name", "Facebook Inc",
  "address", "1 Hacker Way, Menlo Park, CA 94025, USA",
  "created_at", time.Now())

fbhelpdesk := insert("departments", "id",
  "name", "Helpdesk",
  "company_id", facebook,
  "created_at", time.Now())

insert("departments_employees",
  "department_id", fbhelpdesk,
  "employee_id", employee)
```


NOTES

- `facebook` in the above example is the value of the primary key, `companies.id`
- `departments_employees` table has no primary key, so `insert` did not provide a primary key column unlike the other inserts (e.g. `"id"`)
- if repeating `"created_at", time.Now()` seem laborious, then append it to `gofix.Use`

    ``` go
    insert2 := gofix.Use(t, tx, "created_at", time.Now())

    facebook := insert2("companies", "id",
      "name", "Facebook Inc",
      "address", "1 Hacker Way, Menlo Park, CA 94025, USA",
      )

    fbhelpdesk := insert2("departments", "id",
      "name", "Helpdesk",
      "company_id", facebook)
    ```
