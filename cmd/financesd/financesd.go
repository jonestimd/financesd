package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-akka/configuration"
	_ "github.com/go-sql-driver/mysql" // register the driver
	"github.com/graphql-go/handler"
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/graphql"
)

func main() {
	config := configuration.LoadConfig(fmt.Sprintf("%s/.finances/connection.conf", os.Getenv("HOME")))
	driver := strings.ToLower(config.GetString("connection.default.driver"))
	user := config.GetString("connection.default.user")
	password := config.GetString("connection.default.password")
	host := config.GetString("connection.default.host")
	schema := config.GetString("connection.default.schema")
	db, err := gorm.Open(driver, fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, schema))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	db.SingularTable(true)

	graphqlSchema, err := graphql.Schema()
	if err != nil {
		panic(err.Error())
	}
	h := handler.New(&handler.Config{
		Schema:   &graphqlSchema,
		Pretty:   false,
		GraphiQL: true,
	})
	http.Handle("/finances/api/v1/graphql", &graphqlHandler{db: db, handler: h})
	http.ListenAndServe("localhost:8080", nil)
}

type graphqlHandler struct {
	db      *gorm.DB
	handler http.Handler
}

func (h *graphqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// start transaction
	tx := h.db.Begin()
	ctx := context.WithValue(r.Context(), graphql.DbContextKey, tx)
	h.handler.ServeHTTP(w, r.WithContext(ctx))
	// end transaction
	if tx.Error == nil {
		tx.Commit()
	} else {
		tx.Rollback()
	}
}
