// this package is used to perform e2e test
// e2e test will start the server along with
// e2e configs
package e2e

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/e2e"
	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
	"github.com/tanveerprottoy/backend-structure-go/pkg/env"
	"github.com/tanveerprottoy/backend-structure-go/pkg/httpext"
	servere2e "github.com/tanveerprottoy/backend-structure-go/pkg/server/e2e"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const pingTimeout = 10 // in seconds

var (
	db         *sql.DB
	httpClient httpext.Client
	baseURL    string
)

func loadEnv() {
	// as os.Getwd returns the current working directory
	// which is /workers-insights/test/<dir>
	// we need to go back to the root directory
	env.LoadEnv(filepath.Join("..", "..", ".env"))
}

func startPostgresContainer(ctx context.Context) (testcontainers.Container, string, error) {
	pgContainer, err := postgres.Run(ctx,
		"postgres:15.3-alpine",
		postgres.WithDatabase("dummy_ecommerce_db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithInitScripts(filepath.Join("..", "..", "scripts", "db", "init_test.sql")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, "", err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, "", err
	}

	return pgContainer, connStr, nil
}

func initDB(ctx context.Context, connStr string) error {
	var err error
	db, err = sql.Open(constant.DBDriverName, connStr)
	if err != nil {
		return err
	}
	// must ping to check if the connection is successful
	ctx, cancel := context.WithTimeout(ctx, pingTimeout*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping failed with error: %v", err)
	}
	// print the db stats
	stat := db.Stats()
	log.Printf("DB.stats: idle=%d, inUse=%d,  maxOpen=%d", stat.Idle, stat.InUse, stat.MaxOpenConnections)
	return nil
}

func initHTTPClientProvider() {
	httpClient = httpext.NewCustomClient(
		httpext.Config{
			MaxRetries: 5,
			MaxJitter:  10,
			Timeout:    10 * time.Second,
		},
		httpext.WithIdleConnTimeout(50*time.Second),
		httpext.WithMaxIdleConnsPerHost(20),
	)
}

func startServer(chStart, chStop chan int) {
	cfg := e2e.NewConfig(db)
	s := servere2e.NewServer("localhost:"+os.Getenv("PORT"), cfg.Router().Mux)
	s.Start(chStart, chStop)
}

// TestMain is the entry point for the e2e tests
// this function is responsible for setting up the test environment
// and tearing it down after the tests are done
func TestMain(m *testing.M) {
	loadEnv()
	// check if dbtest required to run
	if os.Getenv("E2E_TEST_ENABLED") != "true" {
		log.Println("skipping e2e test")
		os.Exit(0)
	}

	baseURL = os.Getenv("BASE_URL_TEST")

	ctx := context.Background()
	pgContainer, connStr, err := startPostgresContainer(ctx)
	if err != nil {
		log.Printf("err: %v", err)
		os.Exit(1)
	}
	err = initDB(ctx, connStr)
	if err != nil {
		log.Printf("err: %v", err)
		os.Exit(1)
	}
	log.Println("Database initialized, connStr: ", connStr)

	initHTTPClientProvider()

	// start e2e server
	// need to run in a
	// separate goroutine
	// will use 2 channels
	// to signal to start
	// & shutdown the server
	// when tests
	// are run & finished
	chStart := make(chan int)
	chStop := make(chan int)

	go startServer(chStart, chStop)

	// wait for server to start
	<-chStart

	// the below statement will run the tests
	code := m.Run()
	// cleanup code goes here
	// os.Exit doesn't respect defer
	// signal the server to stop
	chStop <- 1
	if err = pgContainer.Terminate(ctx); err != nil {
		log.Printf("err: %v", err)
	}
	log.Println("exiting...")
	// exit with the code from m.Run()
	os.Exit(code)
}
