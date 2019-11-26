package main

import (
	"context"
	"net/http"

	// "os"

	_ "github.com/go-sql-driver/mysql" // because
	"github.com/graphql-go/handler"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/graphql"
)

func main() {
	db, err := gorm.Open("mysql", "user:password@tcp(host)/schema?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	db.SingularTable(true)

	schema, err := graphql.Schema()
	if err != nil {
		panic(err.Error())
	}
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
		RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
			return map[string]interface{}{"db": db}
		},
	})
	http.Handle("/finances/graphql", h)
	http.ListenAndServe("localhost:8080", nil)
}
