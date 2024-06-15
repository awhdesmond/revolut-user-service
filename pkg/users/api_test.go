package users

import (
	"net/http"

	"github.com/stretchr/testify/suite"
	"github.com/upper/db/v4"
)

const (
	apiPrefix = "/hello"
)

type apiTestSuite struct {
	suite.Suite

	pgSess  db.Session
	store   Store
	svc     Service
	handler http.Handler
}

// func (ts *apiTestSuite) SetupSuite() {
// 	logger, _ := common.InitZap("debug")
// 	pgSess, err := common.NewPostgresDBSession(common.TestPgCfg)
// 	if err != nil {
// 		ts.T().Fatalf("got = %v, want = %v", err, nil)
// 	}

// 	store := NewStore(pgSess, logger)
// 	svc := NewService(store, logger)

// 	ts.store = store
// 	ts.svc = svc
// 	// ts.handler, _ MakeHandler(ts.svc, ts.store)
// }

// func (ts *apiTestSuite) TearDownSuite() {
// 	ts.pgSess.SQL().Exec(common.TruncateAllTablesSQL)
// 	time.Sleep(3 * time.Second)
// }
