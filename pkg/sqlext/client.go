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

// Config contains the configuration for the database handle
// it is used as convenience over a number of parameters
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type Option func(*Client)

// WithTLS sets the TLS configuration for the database connection
// rootCert = '/path/to/my/server-ca.pem'
// clientCert = '/path/to/my/client-cert.pem'
// clientKey = '/path/to/my/client-key.pem'
func WithTLS(rootCert, clientCert, clientKey string) Option {
	return func(c *Client) {
		// Configure SSL certificates
		c.dsn = fmt.Sprintf(" sslrootcert=%s sslcert=%s sslkey=%s",
			rootCert, clientCert, clientKey)
	}
}

var (
	instance *Client
	once     sync.Once
)

type Client struct {
	// dsn/DSN = Data Source Name
	dsn string
	db  *sql.DB
}

func GetInstance(opts Config) *Client {
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

func (c *Client) init(opts Config) {
	c.dsn = fmt.Sprintf(
		"host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s",
		opts.Host,
		opts.Port,
		opts.Username,
		opts.Password,
		opts.DBName,
		opts.SSLMode,
	)

	log.Println("dsn: ", c.dsn)

	// Opening a driver typically will not attempt to connect to the database
	// using Must pattern here to ensure that the connection is established
	c.db = must.Must(sql.Open(constant.DBDriverName, c.dsn))

	// Ping the database to verify DSN is valid and the server is accessible
	c.ping(context.Background())

	log.Println("Successfully connected!")

	// set max idle & open connections
	/* d.db.SetMaxIdleConns(maxIdleConns)
	d.db.SetMaxOpenConns(maxOpenConns) */

	stat := c.db.Stats()
	// print the db stats
	log.Printf("DB.stats: idle=%d, inUse=%d,  maxOpen=%d", stat.Idle, stat.InUse, stat.MaxOpenConnections)
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) DB() *sql.DB {
	return c.db
}
