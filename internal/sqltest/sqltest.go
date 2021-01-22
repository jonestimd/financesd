package sqltest

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var whitespace = regexp.MustCompile("[ \n\t]+")
var punctuation = regexp.MustCompile(" +[()] +")

func normalizeSQL(sql string) string {
	normal := whitespace.ReplaceAllLiteralString(sql, " ")
	return punctuation.ReplaceAllStringFunc(normal, func(match string) string {
		return strings.Trim(match, " ")
	})
}

var queryMatcher sqlmock.QueryMatcher = sqlmock.QueryMatcherFunc(func(expectedSQL, actualSQL string) error {
	expectedNormal := normalizeSQL(expectedSQL)
	actualNormal := normalizeSQL(actualSQL)
	if expectedNormal != actualNormal {
		return fmt.Errorf("\n%s\ndoes not match:\n%s", actualNormal, expectedNormal)
	}
	return nil
})

// TestInTx runs a database test in a transaction on a sqlmock connection.
func TestInTx(t *testing.T, test func(mockDB sqlmock.Sqlmock, tx *sql.Tx)) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(queryMatcher))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		panic(err.Error())
	}
	// defer tx.Commit()

	test(mock, tx)
}

// MockRows creates a sqlmock.Rows for the specified columns.
func MockRows(columns ...string) *sqlmock.Rows {
	return sqlmock.NewRows(columns)
}
