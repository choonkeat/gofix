# Experiment to have easier fixtures in Go

``` go
insert := gofix.Use(t, tx)

facebook := insert("companies", "id",
  "name", "Facebook Inc",
  "address", "1 Hacker Way, Menlo Park, CA 94025, USA",
  "created_at", time.Now())

insert("departments", "id",
  "name", "Helpdesk",
  "company_id", facebook,
  "created_at", time.Now())
```

NOTE: `facebook` in the above example is the value of the primary key, `companies.id`
