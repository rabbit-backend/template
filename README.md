# Go SQL Template Engine

This library helps safely build SQL queries in Go by using parameterized queries through Go's powerful templating system.

## Installation

```bash
go get github.com/rabbit-backend/template
```

## Usage

### 1. Define a SQL Template

Create a file named `sql/query.sql.tmpl`:

```sql
SELECT id, username, email
FROM users
WHERE username = {{ .Username | __sql_arg__ }}
AND created_at >= {{ .CreatedAfter | __sql_arg__ }}
LIMIT {{ .Limit | __sql_arg__ }};
```

Create a struct to represent your query parameters:

```go
type QueryParams struct {
	Username     string
	CreatedAfter string
	Limit        int
}
```

### 2. Execute the SQL Query Template

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	engine "github.com/rabbit-backend/template"
	_ "github.com/lib/pq"
)

func main() {
	params := QueryParams{
		Username:     "john_doe",
		CreatedAfter: "2024-01-01",
		Limit:        10,
	}

	sqlEngine := engine.NewEngineWithPlaceHolder(engine.NewPostgresPlaceHolder())

	query, args, err := sqlEngine.Execute("sql/query.sql.tmpl", params)
	if err != nil {
		log.Fatalf("Failed to generate query: %v", err)
	}

	db, err := sql.Open("postgres", "your-connection-string")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	rows, err := db.QueryContext(context.Background(), query, args...)
	if err != nil {
		log.Fatalf("Query execution error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username, email string
		if err := rows.Scan(&id, &username, &email); err != nil {
			log.Fatalf("Row scan error: %v", err)
		}
		fmt.Println(id, username, email)
	}
}
```

---

## Advanced Template Examples

Because this library leverages Go's native `text/template` engine under the hood, you can write expressive dynamic queries using conditionals, loops, and pipeline filters.

### 1. Conditional Logic (`if / else if / else`)

Dynamically include `WHERE` clauses based on incoming parameter values:

```sql
SELECT id, username, email, status, role
FROM users
WHERE 1=1

{{ if .Username }}
  AND username = {{ .Username | __sql_arg__ }}
{{ end }}

{{ if eq .Status "active" }}
  AND is_active = true
{{ else if eq .Status "pending" }}
  AND is_active = false AND verified_at IS NULL
{{ else }}
  AND is_active = false
{{ end }};
```

Corresponding Go struct:

```go
type UserFilterParams struct {
	Username string
	Status   string // "active", "pending", or "inactive"
}
```

### 2. Iteration and Loops (`range`)

Iterate over slices to build dynamic `IN (...)` clauses or bulk `INSERT` statements safely:

#### Dynamic `IN (...)` Clause

```sql
SELECT id, product_name, price
FROM products
WHERE category_id IN (
  {{ range $i, $id := .CategoryIDs }}
    {{ if $i }}, {{ end }}
    {{ $id | __sql_arg__ }}
  {{ end }}
)
AND price <= {{ .MaxPrice | __sql_arg__ }};
```

#### Multi-Row Bulk `INSERT`

```sql
INSERT INTO audit_logs (user_id, action, created_at)
VALUES 
{{ range $i, $log := .Logs }}
  {{ if $i }},{{ end }}
  (
    {{ $log.UserID | __sql_arg__ }}, 
    {{ $log.Action | __sql_arg__ }}, 
    {{ $log.CreatedAt | __sql_arg__ }}
  )
{{ end }};
```

### 3. Pipeline Processing (`pipes`)

Chain Go template functions to transform parameters before binding them to parameterized placeholders:

```sql
SELECT id, title, slug, created_at
FROM posts
WHERE 
  -- Trim whitespace, lower-case, and pass to sql binder
  slug = {{ .Slug | trim | lower | __sql_arg__ }}
  
  -- Wrap search terms in WILDCARD % patterns for LIKE queries
  AND title LIKE {{ .SearchTerm | printf "%%%s%%" | __sql_arg__ }}

ORDER BY created_at DESC
LIMIT {{ .Limit | __sql_arg__ }};
```

---

## Why Parameterized Queries?

Parameterized queries safely handle user inputs, significantly reducing the risk of SQL injection.

### Example of Unsafe Query

```sql
SELECT id FROM users WHERE username = '" + userInput + "';
```

If `userInput` is something malicious, like `admin' OR '1'='1`, it becomes:

```sql
SELECT id FROM users WHERE username = 'admin' OR '1'='1';
```

This condition (`OR '1'='1'`) always evaluates to true, potentially exposing all data.

### Safe Query Using This Library

```sql
SELECT id FROM users WHERE username = $1;
```

With parameters:

```go
["admin' OR '1'='1"]
```

This safely checks the database for the exact string without executing unintended SQL logic.
