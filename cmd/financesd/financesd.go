package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
	if cwd, err := os.Getwd(); err != nil {
		log.Fatal("can't get current directory")
	} else {
		http.Handle("/finances/api/v1/graphql", &graphqlHandler{db: db, handler: h})
		http.Handle("/finances/scripts/", http.StripPrefix("/finances/scripts/", http.FileServer(http.Dir(filepath.Join(cwd, "web", "dist")))))
		http.Handle("/finances/", loadHtml(cwd, config))
		host := config.GetString("listen.host", "localhost")
		port := config.GetInt32("listen.port", 8080)
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil))
	}
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

type staticHtml string

func loadHtml(cwd string, config *configuration.Config) staticHtml {
	htmlTemplate, err := template.New("index.html").ParseFiles(filepath.Join(cwd, "web", "resources", "index.html"))
	if err != nil {
		log.Panicf("Error reading template html: %v", err)
	}
	data := map[string]string{
		"baseUrl": config.GetString("baseUrl", "http://localhost:8080/finances"),
	}
	var buff bytes.Buffer
	if err = htmlTemplate.Execute(&buff, data); err != nil {
		log.Panicf("Error generating html: %v", err)
	}
	return staticHtml(buff.String())
}

func (sf staticHtml) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sendSize := int64(len(sf))
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Content-Length", strconv.FormatInt(sendSize, 10))
	w.WriteHeader(http.StatusOK)
	if r.Method != "HEAD" {
		io.CopyN(w, strings.NewReader(string(sf)), sendSize)
	}
}
