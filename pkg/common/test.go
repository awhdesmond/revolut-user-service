package common

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	TestPgCfg = PostgresSQLConfig{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "postgres",
		Database: "postgres-test",
	}

	TruncateAllTablesSQL = `TRUNCATE TABLE users;`
)

func TestSendReq(req interface{}, path, method string, handler http.Handler) *httptest.ResponseRecorder {
	var httpReq *http.Request
	if req != nil {
		data, _ := json.Marshal(req)
		httpReq = httptest.NewRequest(method, path, bytes.NewReader(data))
	} else {
		httpReq = httptest.NewRequest(method, path, nil)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httpReq)
	return w
}

func TestIsResponseEmptyErr(w *httptest.ResponseRecorder, t *testing.T) {
	TestIsResponseErrorExpected(w, t, "")
}

func TestIsResponseErrorExpected(w *httptest.ResponseRecorder, t *testing.T, wantErr string) {
	res := struct {
		Err string `json:"error"`
	}{}
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatalf("got = % v, want = %v", err, nil)
	}
	if res.Err != wantErr {
		t.Fatalf("got = %v, want = %v", res.Err, wantErr)
	}
}
