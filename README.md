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

	query, args, err := engine.Execute("sql/query.sql.tmpl", params)
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
