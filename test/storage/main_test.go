// this package is used to test the database
// it uses testcontainers to start a postgres container
// this should be run separately from the unit, integration & e2e tests
// as it is meant to test the repository layer with a real database
// this is not really an integration test, nor a unit test
// but it is a test that requires a real database
package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tanveerprottoy/backend-structure-go/pkg/env"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const pingTimeout = 10 // in seconds

var db *sql.DB

func loadEnv() {
	// as os.Getwd returns the current working directory
	// which is ./root/test/<dir>
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
	db, err = sql.Open("pgx", connStr)
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

// TestMain is the entry point for the e2e tests
// this function is responsible for setting up the test environment
// and tearing it down after the tests are done
func TestMain(m *testing.M) {
	loadEnv()

	log.Println(os.Getenv("STORAGE_TEST_ENABLED"))

	// check if storage test is required to run
	if os.Getenv("STORAGE_TEST_ENABLED") != "true" {
		log.Println("skipping storage test")
		os.Exit(0)
	}

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
	// the below statement will run the tests
	code := m.Run()

	// cleanup code goes here
	// os.Exit doesn't respect defer
	if err = pgContainer.Terminate(ctx); err != nil {
		log.Printf("err: %v", err)
	}

	log.Println("exiting...")

	// exit with the code from m.Run()
	os.Exit(code)
}
