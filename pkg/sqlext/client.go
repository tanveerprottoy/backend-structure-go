package sqlext

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/tanveerprottoy/backend-structure-go/pkg/constant"
	"github.com/tanveerprottoy/backend-structure-go/pkg/must"
)

const (
	pingTimeout  = 10 // in seconds
	maxIdleConns = 3
	maxOpenConns = 7
)

type Options struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

var (
	instance *Client
	once     sync.Once
)

type Client struct {
	db *sql.DB
}

func GetInstance(opts Options) *Client {
	once.Do(func() {
		instance = new(Client)
		instance.init(opts)
	})
	return instance
}

// Ping the database to verify DSN is valid and the
// server is accessible. If the ping fails exit the program with an error.
func (c *Client) ping(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, pingTimeout*time.Second)
	defer cancel()
	if err := c.db.PingContext(ctx); err != nil {
		log.Fatalf("ping failed with error: %v", err)
	}
}

func (c *Client) init(opts Options) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s",
		opts.Host,
		opts.Port,
		opts.Username,
		opts.Password,
		opts.DBName,
		opts.SSLMode,
	)
	log.Println("connStr: ", connStr)
	// var err error
	// Opening a driver typically will not attempt to connect to the database.
	c.db = must.Must(sql.Open(constant.DBDriverName, connStr))
	/* if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatalf("db drive open failed with error: %v", err)
	} */
	// Ping the database to verify DSN is valid and the
	// server is accessible
	c.ping(context.Background())
	log.Println("Successfully connected!")
	// set max idle & open connections
	/* d.db.SetMaxIdleConns(maxIdleConns)
	d.db.SetMaxOpenConns(maxOpenConns) */
	// print the db stats
	stat := c.db.Stats()
	log.Printf("DB.stats: idle=%d, inUse=%d,  maxOpen=%d", stat.Idle, stat.InUse, stat.MaxOpenConnections)
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) DB() *sql.DB {
	return c.db
}
