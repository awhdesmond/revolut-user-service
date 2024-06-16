package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
}

func (ts *testSuite) SetupSuite() {
	os.Setenv("REVOLUT_USERS_SVC_POSTGRES_HOST", "localhost")
	os.Setenv("REVOLUT_USERS_SVC_POSTGRES_PORT", "5432")
	os.Setenv("REVOLUT_USERS_SVC_POSTGRES_USERNAME", "postgres")
	os.Setenv("REVOLUT_USERS_SVC_POSTGRES_PASSWORD", "postgres")
	os.Setenv("REVOLUT_USERS_SVC_POSTGRES_DATABASE", "postgres_test")

	go func() {
		main()
	}()
	time.Sleep(2 * time.Second)
}

type ApiTestSuite struct {
	testSuite
}

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

func (ts *ApiTestSuite) Test() {
	ts.testUpsertAndRead()
	ts.testHealth()
}

func (ts *ApiTestSuite) testUpsertAndRead() {
	client := http.Client{}

	// Upsert

	jsonBody := []byte(`{"dateOfBirth": "2020-04-01"}`)
	bodyReader := bytes.NewReader(jsonBody)
	req, err := http.NewRequest(
		http.MethodPut,
		"http://localhost:8080/hello/apple",
		bodyReader,
	)
	if err != nil {
		ts.T().Fatalf("got %v, want = %v", err, nil)
	}

	res, err := client.Do(req)
	if err != nil {
		ts.T().Fatalf("got %v, want = %v", err, nil)
	}
	if res.StatusCode != http.StatusNoContent {
		ts.T().Fatalf("got %v, want = %v", res.StatusCode, http.StatusNoContent)
	}

	// Read

	req, err = http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/hello/apple",
		nil,
	)
	if err != nil {
		ts.T().Fatalf("got %v, want = %v", err, nil)
	}

	res, err = client.Do(req)
	if err != nil {
		ts.T().Fatalf("got %v, want = %v", err, nil)
	}
	if res.StatusCode != http.StatusOK {
		ts.T().Fatalf("got %v, want = %v", res.StatusCode, http.StatusNoContent)
	}

	var bodyResp struct {
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}
	if err := json.NewDecoder(res.Body).Decode(&bodyResp); err != nil {
		ts.T().Fatalf("got %v, want = %v", err, nil)
	}

	if bodyResp.Error != "" {
		ts.T().Fatalf("got %v, want = %v", bodyResp.Error, "")
	}
}

func (ts *ApiTestSuite) testHealth() {
	client := http.Client{}

	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/healthz",
		nil,
	)
	if err != nil {
		ts.T().Fatalf("got %v, want = %v", err, nil)
	}

	res, err := client.Do(req)
	if err != nil {
		ts.T().Fatalf("got %v, want = %v", err, nil)
	}
	if res.StatusCode != http.StatusOK {
		ts.T().Fatalf("got %v, want = %v", res.StatusCode, http.StatusNoContent)
	}
}
