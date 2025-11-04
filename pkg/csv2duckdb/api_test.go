package csv2duckdb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/nalgeon/be"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	dataName      = "gispdata (test).csv"
	dataTableName = "record"
	dataCount     = 131
)

const (
	comTableName     = "event"
	comBrandColumn   = "brand"
	comProductColumn = "product_id"
	comUserIdColumn  = "user_id"
	comHuaweiBrand   = "huawei"
)

const (
	octName                = "2019-Oct (test).csv"
	octCount               = 364
	octHuaweiBrandCount    = 5
	oct1081004856Product   = 1004856
	oct1004856ProductCount = 8
)

const (
	novName             = "2019-Nov (test).csv"
	novCount            = 4880
	novHuaweiBrandCount = 84
)

// getTestDataPath get absolut path fot testdata folder
func getTestDataPath(filename string) string {
	return filepath.Join("..", "..", "testdata", filename)
}

func TestAPI_Import(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	t.Run("Context cancellation", func(t *testing.T) {
		cancel()

		db, err := Import(
			ctx,
			WithTable(getTestDataPath(octName), comTableName),
		)
		be.Err(t, err)
		be.Equal(t, db, nil)

		if db != nil {
			defer func(db *sql.DB) {
				be.Err(t, db.Close(), nil)
			}(db)
		}
	})

	ctx = context.Background()

	t.Run("No file", func(t *testing.T) {
		db, err := Import(ctx)
		be.Err(t, err)
		be.Equal(t, db, nil)
		if db != nil {
			defer func(db *sql.DB) {
				be.Err(t, db.Close(), nil)
			}(db)
		}
	})

	t.Run("Invalid file", func(t *testing.T) {
		db, err := Import(ctx, WithTable("foo.cs", "bar"))
		be.Err(t, err)
		be.Equal(t, db, nil)
		if db != nil {
			defer func(db *sql.DB) {
				be.Err(t, db.Close(), nil)
			}(db)
		}
	})

	t.Run("Progress bar", func(t *testing.T) {
		db, err := Import(
			ctx,
			WithTable(getTestDataPath(octName), comTableName),
			WithProgressBar(false),
		)
		defer func(db *sql.DB) {
			be.Err(t, db.Close(), nil)
		}(db)

		be.Err(t, err, nil)

		// Test import worked without progress bar
		var count int
		err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", comTableName)).Scan(&count)
		be.Err(t, err, nil)
		be.Equal(t, octCount, count)
	})

	t.Run("End count for once table.", func(t *testing.T) {
		db, err := Import(
			ctx,
			WithTable(getTestDataPath(dataName), dataTableName),
		)
		defer func(db *sql.DB) {
			be.Err(t, db.Close(), nil)
		}(db)

		be.Err(t, err, nil)

		var count int
		err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", dataTableName)).Scan(&count)
		be.Err(t, err, nil)

		be.Equal(t, dataCount, count)
	})

	t.Run("Data types are correctly detected", func(t *testing.T) {
		db, err := Import(
			ctx,
			WithTable(getTestDataPath(dataName), dataTableName),
		)
		defer func(db *sql.DB) {
			be.Err(t, db.Close(), nil)
		}(db)

		be.Err(t, err, nil)

		var name string
		var number int
		var date time.Time

		err = db.QueryRow(fmt.Sprintf("SELECT nameoforg, ogrn, docdate FROM %s WHERE registernumber = 10115875", dataTableName)).Scan(&name, &number, &date)
		be.Err(t, err, nil)
		be.Equal(t, "ОБЩЕСТВО С ОГРАНИЧЕННОЙ ОТВЕТСТВЕННОСТЬЮ \"АЙ-ПЛАСТ\"", name)
		be.Equal(t, 1091651001749, number)
		be.Equal(t, time.Date(2023, 5, 24, 0, 0, 0, 0, time.UTC), date) // 2023-05-24
	})

	t.Run("End count for same tables.", func(t *testing.T) {
		db, err := Import(
			ctx,
			WithTables(map[string]string{
				getTestDataPath(octName): comTableName,
				getTestDataPath(novName): comTableName,
			}),
		)
		defer func(db *sql.DB) {
			be.Err(t, db.Close(), nil)
		}(db)

		be.Err(t, err, nil)

		var count int
		err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", comTableName)).Scan(&count)
		be.Err(t, err, nil)
		be.Equal(t, octCount+novCount, count)
	})

	t.Run("Concreate brand count after merge tables.", func(t *testing.T) {
		db, err := Import(
			ctx,
			WithTables(map[string]string{
				getTestDataPath(octName): comTableName,
				getTestDataPath(novName): comTableName,
			}),
		)
		defer func(db *sql.DB) {
			be.Err(t, db.Close(), nil)
		}(db)

		be.Err(t, err, nil)

		var huaweiCount int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE event.%s = '%s'", comTableName,
			comBrandColumn, comHuaweiBrand)
		err = db.QueryRow(query).Scan(&huaweiCount)

		be.Err(t, err, nil)
		be.Equal(t, octHuaweiBrandCount+novHuaweiBrandCount, huaweiCount)
	})

	t.Run("End count for different files with different tables.", func(t *testing.T) {
		db, err := Import(ctx, WithTables(map[string]string{
			getTestDataPath(octName):  comTableName,
			getTestDataPath(dataName): dataTableName,
		}))
		defer func(db *sql.DB) {
			be.Err(t, db.Close(), nil)
		}(db)

		be.Err(t, err, nil)

		var count int
		err = db.QueryRow(fmt.Sprintf("SELECT count(*) FROM %s", comTableName)).Scan(&count)
		be.Err(t, err, nil)
		be.Equal(t, octCount, count)

		err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", dataTableName)).Scan(&count)
		be.Err(t, err, nil)
		be.Equal(t, dataCount, count)
	})

	t.Run("Memory limit", func(t *testing.T) {
		db, err := Import(ctx,
			WithTables(map[string]string{
				getTestDataPath(octName): comTableName,
			}),
			WithMemoryLimit(4))
		defer func(db *sql.DB) {
			be.Err(t, db.Close(), nil)
		}(db)

		be.Err(t, err, nil)

		var count int
		err = db.QueryRow(fmt.Sprintf("SELECT count(*) FROM %s", comTableName)).Scan(&count)
		be.Err(t, err, nil)
		be.Equal(t, octCount, count)
	})

	t.Run("Save as", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "test.duckdb")

		db, err := Import(ctx,
			WithTable(getTestDataPath(octName), comTableName),
			WithSaveAs(tmpFile))
		defer func(db *sql.DB) {
			be.Err(t, db.Close(), nil)
		}(db)

		be.Err(t, err, nil)

		_, err = os.Stat(tmpFile)
		be.Err(t, err, nil)

		// Are you alive?
		var count int
		err = db.QueryRow(fmt.Sprintf("SELECT count(*) FROM %s", comTableName)).Scan(&count)
		be.Err(t, err, nil)
		be.Equal(t, octCount, count)
		be.Err(t, db.Close(), nil)

		current, err := sql.Open("duckdb", tmpFile)
		defer func(db *sql.DB) {
			be.Err(t, db.Close(), nil)
		}(current)
		be.Err(t, err, nil)

		// Are you alive too?
		err = current.QueryRow(fmt.Sprintf("SELECT count(*) FROM %s", comTableName)).Scan(&count)
		be.Err(t, err, nil)
		be.Equal(t, octCount, count)
	})

	t.Run("Column indexes", func(t *testing.T) {
		db, err := Import(ctx,
			WithTable(getTestDataPath(octName), comTableName),
			WithIndexedColumns(map[string][]string{
				comTableName: {comProductColumn, comUserIdColumn},
			}))
		defer func(db *sql.DB) {
			be.Err(t, db.Close(), nil)
		}(db)
		be.Err(t, err, nil)

		// TODO: Make a real check when I get smarter
		var count int
		err = db.QueryRow(
			fmt.Sprintf("SELECT count(*) FROM %s WHERE %s = %d",
				comTableName, comProductColumn, oct1081004856Product)).Scan(&count)
		be.Err(t, err, nil)
		be.Equal(t, oct1004856ProductCount, count)
	})

	t.Run("Insertion order", func(t *testing.T) {
		db, err := Import(ctx,
			WithTable(getTestDataPath(octName), comTableName),
			WithInsertionOrder(true))
		defer func(db *sql.DB) {
			be.Err(t, db.Close(), nil)
		}(db)
		be.Err(t, err, nil)

		var count int
		err = db.QueryRow(fmt.Sprintf("SELECT count(*) FROM %s", comTableName)).Scan(&count)
		be.Err(t, err, nil)
		be.Equal(t, octCount, count)
	})
}
