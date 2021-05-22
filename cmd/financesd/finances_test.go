package main

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MonsantoCo/mocka/v2"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/handler"
	"github.com/jonestimd/financesd/internal/schema"
	"github.com/stretchr/testify/assert"
)

type mockDependencies struct {
	db           *sql.DB
	mockDB       sqlmock.Sqlmock
	staticHTML   *staticHTML
	sqlOpen      *mocka.Stub
	newSchema    *mocka.Stub
	getwd        *mocka.Stub
	httpHandle   *mocka.Stub
	netListen    *mocka.Stub
	serve        *mocka.Stub
	signalNotify func(c chan<- os.Signal, sig ...os.Signal)
	loadHTML     *mocka.Stub
	logAndQuit   func(v ...interface{})
	exitMessage  []interface{}
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
	m.netListen.Restore()
	m.serve.Restore()
	m.loadHTML.Restore()
	signalNotify = m.signalNotify
	logAndQuit = m.logAndQuit
	if verify != nil {
		verify()
	}
}

type mockListener struct {
	closed bool
}

func (l *mockListener) Accept() (net.Conn, error) {
	return nil, &netError{}
}

func (l *mockListener) Close() error {
	l.closed = true
	return nil
}

func (l *mockListener) Addr() net.Addr {
	return nil
}

type netError struct{}

func (ne *netError) Error() string {
	return "fake listener"
}

func (ne *netError) Timeout() bool {
	return false
}

func (ne *netError) Temporary() bool {
	return true
}

func makeMocks(t *testing.T) *mockDependencies {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	staticHTMLValue := &staticHTML{}
	mocks := &mockDependencies{
		db:           db,
		mockDB:       mock,
		staticHTML:   staticHTMLValue,
		sqlOpen:      mocka.Function(t, &sqlOpen, db, nil),
		newSchema:    mocka.Function(t, &newSchema, graphql.Schema{}, nil),
		getwd:        mocka.Function(t, &getwd, "/here", nil),
		httpHandle:   mocka.Function(t, &httpHandle),
		netListen:    mocka.Function(t, &netListen, &mockListener{}, nil),
		serve:        mocka.Function(t, &serve),
		signalNotify: signalNotify,
		loadHTML:     mocka.Function(t, &loadHTML, staticHTMLValue),
		logAndQuit:   logAndQuit,
	}
	signalNotify = func(c chan<- os.Signal, sig ...os.Signal) {
		c <- os.Interrupt
	}
	logAndQuit = func(v ...interface{}) {
		mocks.exitMessage = v
		panic("log.Fatal")
	}
	return mocks
}

func Test_main_connectsToDatabaseAndStartsServer(t *testing.T) {
	mocks := makeMocks(t)
	mocks.mockDB.ExpectPing()
	defer mocks.restore(t, "", nil)
	os.Args = os.Args[0:1]

	main()

	assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
	assert.Equal(t, 3, mocks.httpHandle.CallCount())
	assert.Equal(t, []interface{}{"/finances/api/v1/graphql"}, mocks.httpHandle.GetCall(0).Arguments()[:1])
	assert.Equal(t, mocks.db, mocks.httpHandle.GetCall(0).Arguments()[1].(*graphqlHandler).db)
	assert.Equal(t, []interface{}{"/finances/scripts/"}, mocks.httpHandle.GetCall(1).Arguments()[:1])
	assert.Equal(t, []interface{}{"/finances/", mocks.staticHTML}, mocks.httpHandle.GetCall(2).Arguments())
	assert.Equal(t, 1, mocks.netListen.CallCount())
	assert.Equal(t, []interface{}{"tcp", "localhost:8080"}, mocks.netListen.GetCall(0).Arguments())
	assert.Nil(t, mocks.exitMessage)
}

func Test_main_quitsIfDbConnectionFails(t *testing.T) {
	mocks := makeMocks(t)
	expectedErr := errors.New("config error")
	mocks.sqlOpen.OnFirstCall().Return(nil, expectedErr)
	defer mocks.restore(t, "log.Fatal", func() {
		assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
		assert.Equal(t, []interface{}{expectedErr}, mocks.exitMessage)
	})
	os.Args = os.Args[0:1]

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
	os.Args = os.Args[0:1]

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
	os.Args = os.Args[0:1]

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
	os.Args = os.Args[0:1]

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
	os.Args = os.Args[0:1]

	main()
}

type mockGraphql struct {
	user     interface{}
	setError bool
	panic    bool
}

func (h *mockGraphql) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.user = r.Context().Value(schema.UserKey)
	if h.setError {
		hasError := r.Context().Value(hasErrorKey).(*bool)
		*hasError = true
	} else if h.panic {
		panic("graphql error")
	}
}

func Test_ServeHTTP_requiresUser(t *testing.T) {
	tests := []struct {
		name        string
		requestUser string
		defaultUser string
		status      int
	}{
		{"no user", "", "", http.StatusBadRequest},
		{"request user", "somebody", "nobody", http.StatusOK},
		{"default user", "", "somebody", http.StatusOK},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mocks := makeMocks(t)
			defer mocks.restore(t, "l", nil)
			if test.status == http.StatusOK {
				mocks.mockDB.ExpectBegin()
				mocks.mockDB.ExpectCommit()
			}
			gqlHandler := &mockGraphql{}
			handler := &graphqlHandler{db: mocks.db, defaultUser: test.defaultUser, handler: gqlHandler}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, &http.Request{Header: http.Header{"X-User": {test.requestUser}}})

			assert.Equal(t, test.status, w.Result().StatusCode)
			assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
			if test.status == http.StatusOK {
				assert.Equal(t, "somebody", gqlHandler.user)
			}
		})
	}
}

func Test_ServeHTTP_returnsDatabaseError(t *testing.T) {
	mocks := makeMocks(t)
	defer mocks.restore(t, "l", nil)
	handler := &graphqlHandler{db: mocks.db, defaultUser: "somebody", handler: &mockGraphql{}}
	mocks.mockDB.ExpectBegin().WillReturnError(errors.New("connection error"))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, &http.Request{})

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
}

func Test_ServeHTTP_commitsOnSuccess(t *testing.T) {
	mocks := makeMocks(t)
	defer mocks.restore(t, "", nil)
	handler := &graphqlHandler{db: mocks.db, defaultUser: "somebody", handler: &mockGraphql{}}
	mocks.mockDB.ExpectBegin()
	mocks.mockDB.ExpectCommit()
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, &http.Request{})

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
}

func Test_ServeHTTP_returnsErrorIfCommitFails(t *testing.T) {
	mocks := makeMocks(t)
	defer mocks.restore(t, "", nil)
	handler := &graphqlHandler{db: mocks.db, defaultUser: "somebody", handler: &mockGraphql{}}
	mocks.mockDB.ExpectBegin()
	mocks.mockDB.ExpectCommit().WillReturnError(errors.New("commit error"))
	h := httptest.NewRecorder()

	handler.ServeHTTP(h, &http.Request{})

	assert.Equal(t, http.StatusInternalServerError, h.Result().StatusCode)
	assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
}

func Test_ServeHTTP_rollsbackOnError(t *testing.T) {
	mocks := makeMocks(t)
	defer mocks.restore(t, "", nil)
	handler := &graphqlHandler{db: mocks.db, defaultUser: "somebody", handler: &mockGraphql{setError: true}}
	mocks.mockDB.ExpectBegin()
	mocks.mockDB.ExpectRollback()
	h := httptest.NewRecorder()

	handler.ServeHTTP(h, &http.Request{})

	assert.Equal(t, http.StatusOK, h.Result().StatusCode, "doesn't change status")
	assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
}

func Test_ServeHTTP_rollsbackOnPanic(t *testing.T) {
	mocks := makeMocks(t)
	defer mocks.restore(t, "graphql error", nil)
	handler := &graphqlHandler{db: mocks.db, defaultUser: "somebody", handler: &mockGraphql{panic: true}}
	mocks.mockDB.ExpectBegin()
	mocks.mockDB.ExpectRollback()
	h := httptest.NewRecorder()

	handler.ServeHTTP(h, &http.Request{})

	assert.Equal(t, http.StatusOK, h.Result().StatusCode, "doesn't change status")
	assert.Nil(t, mocks.mockDB.ExpectationsWereMet())
}
