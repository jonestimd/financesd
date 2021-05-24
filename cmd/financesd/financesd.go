package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/go-akka/configuration"
	_ "github.com/go-sql-driver/mysql" // register the driver
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/jonestimd/financesd/internal/schema"
	"github.com/jonestimd/financesd/internal/server"
)

var httpHandle = http.Handle
var netListen = net.Listen
var signalNotify = signal.Notify
var logAndQuit = log.Fatal
var sqlOpen = sql.Open
var newSchema = schema.New
var getwd = os.Getwd
var newHandler = handler.New

func main() {
	configPath := fmt.Sprintf("%s/.finances/connection.conf", os.Getenv("HOME"))
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	config := configuration.LoadConfig(configPath)
	driver := strings.ToLower(config.GetString("connection.default.driver"))
	user := config.GetString("connection.default.user")
	password := config.GetString("connection.default.password")
	host := config.GetString("connection.default.host")
	schema := config.GetString("connection.default.schema")
	db, err := sqlOpen(driver, fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, schema))
	if err != nil {
		logAndQuit(err)
	}
	defer db.Close()
	db.SetConnMaxLifetime(config.GetTimeDuration("connection.maxLifetime", 0))
	db.SetMaxIdleConns(int(config.GetInt32("connection.maxIdleConnections", 10)))
	db.SetMaxOpenConns(int(config.GetInt32("connection.maxOpenConnections", 10)))
	if err := db.Ping(); err != nil {
		logAndQuit(err)
	}

	graphqlSchema, err := newSchema()
	if err != nil {
		logAndQuit(err)
	}
	gqlHandler := newHandler(&handler.Config{
		Schema:           &graphqlSchema,
		Pretty:           false,
		GraphiQL:         true,
		ResultCallbackFn: resultCallback,
	})
	if cwd, err := getwd(); err != nil {
		logAndQuit("can't get current directory")
	} else {
		network, address := getListenConfig(config.GetConfig("listen"))
		httpHandle("/finances/api/v1/graphql", &graphqlHandler{db: db, defaultUser: config.GetString("oauth.user"), handler: gqlHandler})
		httpHandle("/finances/scripts/", http.StripPrefix("/finances/scripts/", http.FileServer(http.Dir(filepath.Join(cwd, "web", "dist")))))
		httpHandle("/finances/", newIndexHandler(cwd, network, address))
		umask, err := strconv.ParseInt(config.GetString("listen.umask", "0117"), 8, 32)
		if err != nil {
			umask = 79
		}
		syscall.Umask(int(umask))
		log.Printf("Listening at %s:%s\n", network, address)
		listener, err := netListen(network, address)
		if err != nil {
			log.Fatalf("Error listening: %v", err)
		}
		defer listener.Close()
		idSupplier := server.NewIDSupplier()
		go func() {
			router := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestID := idSupplier.Next()
				snoop := httpsnoop.CaptureMetrics(http.DefaultServeMux, w, r.WithContext(context.WithValue(r.Context(), requestIdKey, requestID)))
				log.Printf("[%s] %d %-5s %s %s %d %v %s\n", requestID, snoop.Code, r.Method, r.URL.Path, r.Host, snoop.Written,
					snoop.Duration.Truncate(time.Millisecond), r.UserAgent())
			})
			serve(listener, router)
		}()

		sigc := make(chan os.Signal, 1)
		signalNotify(sigc, os.Interrupt, syscall.SIGTERM)
		s := <-sigc
		log.Print("Got signal: ", s)
	}
}

var serve = func(listener net.Listener, router http.HandlerFunc) {
	server := &http.Server{Handler: router}
	log.Fatal(server.Serve(listener))
}

func getListenConfig(config *configuration.Config) (network string, address string) {
	network = config.GetString("network", "tcp")
	address = config.GetString("address", "localhost:8080")
	return
}

type graphqlHandler struct {
	db          *sql.DB
	defaultUser string
	handler     http.Handler
}

type reqContextKey string

const requestIdKey = reqContextKey("requestID")
const hasErrorKey = reqContextKey("hasError")

var compressSpaces = regexp.MustCompile(`\s+`)

func resultCallback(ctx context.Context, params *graphql.Params, result *graphql.Result, responseBody []byte) {
	requestID := ctx.Value(requestIdKey)
	log.Printf("[%s] query: %s", requestID, compressSpaces.ReplaceAllString(params.RequestString, " "))
	log.Printf("[%s] variables: %v", requestID, params.VariableValues)
	hasError := ctx.Value(hasErrorKey).(*bool)
	*hasError = result.HasErrors()
}

func (h *graphqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var hasError bool
	ctx := context.WithValue(r.Context(), hasErrorKey, &hasError)
	if user := r.Header.Get("X-User"); user != "" {
		ctx = context.WithValue(ctx, schema.UserKey, user)
	} else if h.defaultUser != "" {
		ctx = context.WithValue(ctx, schema.UserKey, h.defaultUser)
	} else {
		http.Error(w, "Unknown user", http.StatusBadRequest)
		return
	}
	// start transaction
	tx, err := h.db.Begin()
	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	ctx = context.WithValue(ctx, schema.DbContextKey, tx)
	h.handler.ServeHTTP(w, r.WithContext(ctx))
	// end transaction
	requestID := ctx.Value(requestIdKey)
	if hasError {
		log.Printf("[%s] Rolling back", requestID)
		if err := tx.Rollback(); err != nil {
			log.Printf("[%s] Rollback failed: %v", requestID, err)
		}
	} else if err := tx.Commit(); err != nil {
		log.Printf("[%s] Commit failed: %v", requestID, err)
		http.Error(w, fmt.Sprintf("Commit failed: %v", err), http.StatusInternalServerError)
	}
}

type staticHTML struct {
	modTime time.Time
	size    int64
	content string
}

type staticHandler struct {
	baseUrl  string
	template *template.Template
	renders  map[string]*staticHTML
}

var newIndexHandler = func(cwd string, network string, address string) *staticHandler {
	file := filepath.Join(cwd, "web", "resources", "index.html")
	htmlTemplate, err := template.New("index.html").ParseFiles(file)
	if err != nil {
		log.Panicf("Error reading template html: %v", err)
	}
	var baseUrl string
	if network == "tcp" {
		baseUrl = "http://" + address + "/finances/"
	} else {
		baseUrl = "unix://" + address + "/finances"
	}
	return &staticHandler{baseUrl: baseUrl, template: htmlTemplate, renders: make(map[string]*staticHTML)}
}

var unixEpochTime = time.Unix(0, 0)

func (st *staticHandler) getHTML(r *http.Request) *staticHTML {
	baseUrl := r.Header.Get("X-Forwarded-For")
	if baseUrl == "" {
		baseUrl = st.baseUrl
	}
	if baseUrl[len(baseUrl)-1] == '/' {
		baseUrl = baseUrl[:len(baseUrl)-1]
	}
	var html *staticHTML
	var ok bool
	if html, ok = st.renders[baseUrl]; !ok {
		data := map[string]string{"baseUrl": baseUrl}
		var buff bytes.Buffer
		if err := st.template.Execute(&buff, data); err != nil {
			log.Panicf("Error generating html: %v", err)
		}
		content := buff.String()
		html = &staticHTML{content: content, modTime: time.Now(), size: int64(len(content))}
		st.renders[baseUrl] = html
	}
	return html
}

func (st *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" && r.Method != "OPTIONS" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Allow", "GET, OPTIONS, HEAD")
	}
	html := st.getHTML(r)
	w.Header().Set("Last-Modified", html.modTime.UTC().Format(http.TimeFormat))
	ims := r.Header.Get("If-Modified-Since")
	if ims != "" && !html.modTime.IsZero() && html.modTime != unixEpochTime {
		if t, err := http.ParseTime(ims); err == nil {
			// The Date-Modified header truncates sub-second precision, so
			// use mtime < t+1s instead of mtime <= t to check for unmodified.
			if !html.modTime.Before(t.Add(1 * time.Second)) {
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
	w.Header().Set("Content-Length", strconv.FormatInt(html.size, 10))
	w.WriteHeader(http.StatusOK)
	if r.Method == "GET" {
		io.CopyN(w, strings.NewReader(html.content), html.size)
	}
}
