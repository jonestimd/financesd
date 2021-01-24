package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/go-akka/configuration"
	_ "github.com/go-sql-driver/mysql" // register the driver
	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/jonestimd/financesd/internal/graphql"
)

func main() {
	config := configuration.LoadConfig(fmt.Sprintf("%s/.finances/connection.conf", os.Getenv("HOME")))
	driver := strings.ToLower(config.GetString("connection.default.driver"))
	user := config.GetString("connection.default.user")
	password := config.GetString("connection.default.password")
	host := config.GetString("connection.default.host")
	schema := config.GetString("connection.default.schema")
	db, err := sql.Open(driver, fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, schema))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(config.GetTimeDuration("connection.maxLifetime", 0))
	db.SetMaxIdleConns(int(config.GetInt32("connection.maxIdleConnections", 10)))
	db.SetMaxOpenConns(int(config.GetInt32("connection.maxOpenConnections", 10)))
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	graphqlSchema, err := graphql.Schema()
	if err != nil {
		log.Fatal(err)
	}
	gqlHandler := handler.New(&handler.Config{
		Schema:           &graphqlSchema,
		Pretty:           false,
		GraphiQL:         true,
		ResultCallbackFn: resultCallback,
		RootObjectFn: handler.RootObjectFn(func(ctx context.Context, r *http.Request) map[string]interface{} {
			return make(map[string]interface{}) // TODO use struct
		}),
	})
	if cwd, err := os.Getwd(); err != nil {
		log.Fatal("can't get current directory")
	} else {
		http.Handle("/finances/api/v1/graphql", &graphqlHandler{db: db, handler: gqlHandler})
		http.Handle("/finances/scripts/", http.StripPrefix("/finances/scripts/", http.FileServer(http.Dir(filepath.Join(cwd, "web", "dist")))))
		http.Handle("/finances/", loadHTML(cwd, config))
		router := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			snoop := httpsnoop.CaptureMetrics(http.DefaultServeMux, w, r)
			log.Printf("%d %-5s %s %s %d %v %s\n", snoop.Code, r.Method, r.URL.Path, r.Host, snoop.Written,
				snoop.Duration.Truncate(time.Millisecond), r.UserAgent())
		})
		host := config.GetString("listen.host", "localhost")
		port := config.GetInt32("listen.port", 8080)
		log.Printf("Listening at %s:%d/finances\n", host, port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), router))
	}
}

type graphqlHandler struct {
	db      *sql.DB
	handler http.Handler
}

type reqContextKey string

const hasErrorKey = reqContextKey("hasError")

func resultCallback(ctx context.Context, params *gql.Params, result *gql.Result, responseBody []byte) {
	hasError := ctx.Value(hasErrorKey).(*bool)
	*hasError = result.HasErrors()
}

func (h *graphqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// start transaction
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
	} else {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				panic(r)
			}
		}()
		var hasError bool
		ctx := context.WithValue(r.Context(), graphql.DbContextKey, tx)
		ctx = context.WithValue(ctx, hasErrorKey, &hasError)
		h.handler.ServeHTTP(w, r.WithContext(ctx))
		// end transaction
		if hasError {
			log.Println("Rolling back")
			if err := tx.Rollback(); err != nil {
				log.Printf("Rollback failed: %v\n", err)
			}
		} else if err := tx.Commit(); err != nil {
			log.Printf("Commit failed: %v\n", err)
			http.Error(w, fmt.Sprintf("Commit failed: %v", err), http.StatusInternalServerError)
		}
	}
}

type staticHTML struct {
	modTime time.Time
	size    int64
	content string
}

func loadHTML(cwd string, config *configuration.Config) *staticHTML {
	file := filepath.Join(cwd, "web", "resources", "index.html")
	stat, err := os.Stat(file)
	if err != nil {
		log.Panicf("Can't stat file: %s, %v", file, err)
	}
	htmlTemplate, err := template.New("index.html").ParseFiles(file)
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
	content := buff.String()
	return &staticHTML{content: content, size: int64(len(content)), modTime: stat.ModTime()}
}

var unixEpochTime = time.Unix(0, 0)

func (st *staticHTML) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" && r.Method != "OPTIONS" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Allow", "GET, OPTIONS, HEAD")
	}
	w.Header().Set("Last-Modified", st.modTime.UTC().Format(http.TimeFormat))
	ims := r.Header.Get("If-Modified-Since")
	if ims != "" && !st.modTime.IsZero() && st.modTime != unixEpochTime {
		if t, err := http.ParseTime(ims); err == nil {
			// The Date-Modified header truncates sub-second precision, so
			// use mtime < t+1s instead of mtime <= t to check for unmodified.
			if !st.modTime.Before(t.Add(1 * time.Second)) {
				h := w.Header()
				delete(h, "Content-Type")
				delete(h, "Content-Length")
				if h.Get("Etag") != "" {
					delete(h, "Last-Modified")
				}
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Content-Length", strconv.FormatInt(st.size, 10))
	w.WriteHeader(http.StatusOK)
	if r.Method == "GET" {
		io.CopyN(w, strings.NewReader(st.content), st.size)
	}
}
