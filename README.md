# csv2duckdb
[![Go Reference](https://pkg.go.dev/badge/github.com/ev-kotov/csv2duckdb.svg)](https://pkg.go.dev/github.com/ev-kotov/csv2duckdb)
[![Go Report Card](https://goreportcard.com/badge/github.com/ev-kotov/csv2duckdb)](https://goreportcard.com/report/github.com/ev-kotov/csv2duckdb)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A simple and minimalistic **Go** library for importing large CSV files into [DuckDB](https://github.com/duckdb/duckdb).

## Highlights
- special **optimized** for working with **big** data;
- **minimalistic API**, just one func with seven optional settings;
- two operating modes:
  - "in-memory"
  - saved to a local database to return later;
- zero hassle.

## Installation
```bash
go get github.com/ev-kotov/csv2duckdb
```

## Usage
### Import once table "in-memory"
- file - `foo.csv`
- table - `bar`

```go
package main

import (
  "context"
  "fmt"
  "github.com/ev-kotov/csv2duckdb/pkg/csv2duckdb"
)

func main() {
  ctx := context.Background()
  db, _ := csv2duckdb.Import(ctx,csv2duckdb.WithTable("foo.csv", "bar"))
  
  var count int
  _ = db.QueryRow("SELECT COUNT(*) FROM BAR").Scan(&count)
  fmt.Println(count)
}
```

### Import same tables
- with various files: `foo.csv`, `bar.csv`
- with various tables: `baz`, `qux`

```go
package main

import (
  "context"
  "fmt"
  "github.com/ev-kotov/csv2duckdb/pkg/csv2duckdb"
)

func main() {
  ctx := context.Background()
  db, _ := csv2duckdb.Import(ctx,
    csv2duckdb.WithTables(map[string]string{
      "foo.csv": "baz",
      "bar.csv": "qux"}))

  var count int
  _ = db.QueryRow("SELECT COUNT(*) FROM(SELECT * FROM BAZ UNION ALL SELECT * FROM QUX) AS result").Scan(&count)
  fmt.Println(count)
}
```
- with various files: `foo.csv`, `bar.csv`
- with identical table: `baz`

```go
package main

import (
  "context"
  "fmt"
  "github.com/ev-kotov/csv2duckdb/pkg/csv2duckdb"
)

func main() {
  ctx := context.Background()
  db, _ := csv2duckdb.Import(ctx,
    csv2duckdb.WithTables(map[string]string{
      "foo.csv": "baz",
      "bar.csv": "baz"})) // automatic merging

  var count int
  _ = db.QueryRow("SELECT COUNT(*) FROM BAZ").Scan(&count)
  fmt.Println(count)
}
```

### local save database

```go
package main

import (
  "context"
  "github.com/ev-kotov/csv2duckdb/pkg/csv2duckdb"
  "os"
)

func main() {
  ctx := context.Background()
  _, _ = csv2duckdb.Import(ctx,
    csv2duckdb.WithTable("foo.csv", "bar"),
    csv2duckdb.WithSaveAs(os.TempDir()))
}
```

### Column indexes
- files - `foo.csv`, `bar.csv`
- table - `baz`, `qux`
- column indexes: `quux`, `corge`, `grault`, `garply`
```go
package main

import (
	"context"
	"fmt"
	"github.com/ev-kotov/csv2duckdb/pkg/csv2duckdb"
)

func main() {
	ctx := context.Background()
	db, _ := csv2duckdb.Import(ctx,
		csv2duckdb.WithTables(map[string]string{
			"foo.csv": "baz",
			"bar.csv": "qux"}),
		csv2duckdb.WithIndexedColumns(map[string][]string{
			"baz": {"quux", "corge"},
			"qux": {"grault", "garply"}}))

	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM BAR").Scan(&count)
	fmt.Println(count)
}
```
## Contributing
Bug fixes are welcome. 
For anything other than bug fixes, please open an issue first to discuss your proposed changes. 
The package has a very limited scope, so it's important to discuss any new features before implementing them.

Make sure to add or update tests as needed.

