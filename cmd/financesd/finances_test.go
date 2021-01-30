package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/handler"
	"github.com/stretchr/testify/assert"
)

type mockDependencies struct {
	db                 *sql.DB
	mockDB             sqlmock.Sqlmock
	staticHTML         *staticHTML
	sqlOpen            *mocka.Stub
	newSchema          *mocka.Stub
	getwd              *mocka.Stub
	httpHandle         *mocka.Stub
	httpListenAndServe *mocka.Stub
	loadHTML           *mocka.Stub
	logAndQuit         func(v ...interface{})
	exitMessage        []interface{}
}

func (m *mockDependencies) restore(t *testing.T, panicMessage string, verify func()) {
	if r := recover(); r != nil {
		assert.Equal(t, panicMessage, r)
	}
	m.db.Close()
	m.sqlOpen.Restore()
	m.newSchema.Restore()
	m.getwd.Restore()
	m.httpHandle.Restore()
	m.httpListenAndServe.Restore()
	m.loadHTML.Restore()
	logAndQuit = m.logAndQuit
	if verify != nil {
		verify()
	}
}

func makeMocks(t *testing.T) *mockDependencies {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	staticHTMLValue := &staticHTML{}
	mocks := &mockDependencies{
		db:                 db,
		mockDB:             mock,
		staticHTML:         staticHTMLValue,
		sqlOpen:            mocka.Function(t, &sqlOpen, db, nil),
		newSchema:          mocka.Function(t, &newSchema, graphql.Schema{}, nil),
		getwd:              mocka.Function(t, &getwd, "/here", nil),
		httpHandle:         mocka.Function(t, &httpHandle),
		httpListenAndServe: mocka.Function(t, &httpListenAndServe, errors.New("stopped server")),
		loadHTML:           mocka.Function(t, &loadHTML, staticHTMLValue),
	}
	mocks.logAndQuit = logAndQuit
	logAndQuit = func(v ...interface{}) {
		mocks.exitMessage = v
		panic("log.Fatal")
	}
	return mocks
}

func Test_main_connectsToDatabaseAndStartsServer(t *testing.T) {
	mocks := makeMocks(t)
	mocks.mockDB.ExpectPing()
	defer mocks.restore(t, "log.Fatal", func() {
		assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
		assert.Equal(t, 3, mocks.httpHandle.CallCount())
		assert.Equal(t, []interface{}{"/finances/api/v1/graphql"}, mocks.httpHandle.GetCall(0).Arguments()[:1])
		assert.Equal(t, mocks.db, mocks.httpHandle.GetCall(0).Arguments()[1].(*graphqlHandler).db)
		assert.Equal(t, []interface{}{"/finances/scripts/"}, mocks.httpHandle.GetCall(1).Arguments()[:1])
		assert.Equal(t, []interface{}{"/finances/", mocks.staticHTML}, mocks.httpHandle.GetCall(2).Arguments())
		assert.Equal(t, 1, mocks.httpListenAndServe.CallCount())
		assert.Equal(t, []interface{}{"localhost:8080"}, mocks.httpListenAndServe.GetCall(0).Arguments()[:1])
		assert.Equal(t, []interface{}{errors.New("stopped server")}, mocks.exitMessage)
	})

	main()

	assert.Fail(t, "expected log.Fatal")
}

func Test_main_quitsIfDbConnectionFails(t *testing.T) {
	mocks := makeMocks(t)
	expectedErr := errors.New("config error")
	mocks.sqlOpen.OnFirstCall().Return(nil, expectedErr)
	defer mocks.restore(t, "log.Fatal", func() {
		assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
		assert.Equal(t, []interface{}{expectedErr}, mocks.exitMessage)
	})

	main()

	assert.Fail(t, "expected log.Fatal")
}

func Test_main_quitsIfDbPingFails(t *testing.T) {
	mocks := makeMocks(t)
	expectedErr := errors.New("connection error")
	mocks.mockDB.ExpectPing().WillReturnError(expectedErr)
	defer mocks.restore(t, "log.Fatal", func() {
		assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
		assert.Equal(t, []interface{}{expectedErr}, mocks.exitMessage)
	})

	main()

	assert.Fail(t, "expected log.Fatal")
}

func Test_main_quitsOnSchemaError(t *testing.T) {
	mocks := makeMocks(t)
	mocks.mockDB.ExpectPing()
	expectedErr := errors.New("schema error")
	mocks.newSchema.OnFirstCall().Return(graphql.Schema{}, expectedErr)
	defer mocks.restore(t, "log.Fatal", func() {
		assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
		assert.Equal(t, []interface{}{expectedErr}, mocks.exitMessage)
	})

	main()

	assert.Fail(t, "expected log.Fatal")
}

func Test_main_quitsIfCurrentDirUnavailable(t *testing.T) {
	mocks := makeMocks(t)
	mocks.mockDB.ExpectPing()
	expectedErr := errors.New("getwd error")
	mocks.getwd.OnFirstCall().Return("", expectedErr)
	defer mocks.restore(t, "log.Fatal", func() {
		assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
		assert.Equal(t, []interface{}{"can't get current directory"}, mocks.exitMessage)
	})

	main()

	assert.Fail(t, "expected log.Fatal")
}

func Test_resultCallback(t *testing.T) {
	mockNewHandler := mocka.Function(t, &newHandler, &handler.Handler{})
	defer mockNewHandler.Restore()
	mocks := makeMocks(t)
	mocks.mockDB.ExpectPing()
	defer mocks.restore(t, "log.Fatal", func() {
		callback := mockNewHandler.GetCall(0).Arguments()[0].(*handler.Config).ResultCallbackFn
		var hasError bool
		ctx := context.WithValue(context.TODO(), hasErrorKey, &hasError)

		callback(ctx, nil, &graphql.Result{}, nil)
		assert.Equal(t, false, hasError)

		callback(ctx, nil, &graphql.Result{Errors: []gqlerrors.FormattedError{{}}}, nil)
		assert.Equal(t, true, hasError)
	})
	main()
}

type mockResponse struct {
	header map[string][]string
	status int
}

func (r *mockResponse) Header() http.Header {
	if r.header == nil {
		r.header = make(map[string][]string)
	}
	return r.header
}

func (r *mockResponse) Write([]byte) (int, error) {
	return 0, nil
}

func (r *mockResponse) WriteHeader(statusCode int) {
	r.status = statusCode
}

type mockGraphql struct {
	setError bool
	panic    bool
}

func (h *mockGraphql) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.setError {
		hasError := r.Context().Value(hasErrorKey).(*bool)
		*hasError = true
	} else if h.panic {
		panic("graphql error")
	}
}

func Test_ServeHTTP_returnsDatabaseError(t *testing.T) {
	mocks := makeMocks(t)
	defer mocks.restore(t, "l", nil)
	handler := &graphqlHandler{db: mocks.db, handler: &mockGraphql{}}
	mocks.mockDB.ExpectBegin().WillReturnError(errors.New("connection error"))
	res := &mockResponse{}

	handler.ServeHTTP(res, nil)

	assert.Equal(t, http.StatusInternalServerError, res.status)
	assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
}

func Test_ServeHTTP_commitsOnSuccess(t *testing.T) {
	mocks := makeMocks(t)
	defer mocks.restore(t, "", nil)
	handler := &graphqlHandler{db: mocks.db, handler: &mockGraphql{}}
	mocks.mockDB.ExpectBegin()
	mocks.mockDB.ExpectCommit()
	res := &mockResponse{}

	handler.ServeHTTP(res, &http.Request{})

	assert.Equal(t, 0, res.status)
	assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
}

func Test_ServeHTTP_returnsErrorIfCommitFails(t *testing.T) {
	mocks := makeMocks(t)
	defer mocks.restore(t, "", nil)
	handler := &graphqlHandler{db: mocks.db, handler: &mockGraphql{}}
	mocks.mockDB.ExpectBegin()
	mocks.mockDB.ExpectCommit().WillReturnError(errors.New("commit error"))
	res := &mockResponse{}

	handler.ServeHTTP(res, &http.Request{})

	assert.Equal(t, http.StatusInternalServerError, res.status)
	assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
}

func Test_ServeHTTP_rollsbackOnError(t *testing.T) {
	mocks := makeMocks(t)
	defer mocks.restore(t, "", nil)
	handler := &graphqlHandler{db: mocks.db, handler: &mockGraphql{setError: true}}
	mocks.mockDB.ExpectBegin()
	mocks.mockDB.ExpectRollback()
	res := &mockResponse{}

	handler.ServeHTTP(res, &http.Request{})

	assert.Equal(t, 0, res.status, "doesn't change status")
	assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
}

func Test_ServeHTTP_rollsbackOnPanic(t *testing.T) {
	mocks := makeMocks(t)
	defer mocks.restore(t, "graphql error", nil)
	handler := &graphqlHandler{db: mocks.db, handler: &mockGraphql{panic: true}}
	mocks.mockDB.ExpectBegin()
	mocks.mockDB.ExpectRollback()
	res := &mockResponse{}

	handler.ServeHTTP(res, &http.Request{})

	assert.Equal(t, 0, res.status, "doesn't change status")
	assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
}
