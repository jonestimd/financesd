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
	})
	http.Handle("/finances/graphql", &graphqlHandler{db: db, handler: h})
	http.ListenAndServe("localhost:8080", nil)
}

type graphqlHandler struct {
	db      *gorm.DB
	handler http.Handler
}

func (h *graphqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO start transaction
	h.handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), graphql.DbContextKey, h.db)))
	// TODO end transaction
}
