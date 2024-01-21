package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/dot96gal/go-sqlc-sample/internal/sqlc"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *sql.DB

func TestMain(m *testing.M) {
	mysqlDatabase := os.Getenv("TEST_MYSQL_DATABASE")
	mysqlRootPass := os.Getenv("TEST_MYSQL_ROOT_PASSWORD")
	mysqlUser := os.Getenv("TEST_MYSQL_USER")
	mysqlPass := os.Getenv("TEST_MYSQL_PASSWORD")
	mysqlHost := os.Getenv("TEST_MYSQL_HOST")
	mysqlPort := os.Getenv("TEST_MYSQL_TCP_PORT")

	pwd, _ := os.Getwd()
	schemaFiles := fmt.Sprintf("file://%s/db/migrations", pwd)

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	pool.MaxWait = 30 * time.Second

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	runOptions := &dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.3.0",
		Env: []string{
			fmt.Sprintf("MYSQL_DATABASE=%s", mysqlDatabase),
			fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", mysqlRootPass),
			fmt.Sprintf("MYSQL_USER=%s", mysqlUser),
			fmt.Sprintf("MYSQL_PASSWORD=%s", mysqlPass),
			fmt.Sprintf("MYSQL_HOST=%s", mysqlHost),
			fmt.Sprintf("MYSQL_TCP_PORT=%s", mysqlPort),
		},
	}

	resource, err := pool.RunWithOptions(
		runOptions,
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	port := resource.GetPort(fmt.Sprintf("%s/tcp", mysqlPort))
	dataSource := fmt.Sprintf("%s:%s@(%s:%s)/%s", mysqlUser, mysqlPass, mysqlHost, port, mysqlDatabase)

	if err := pool.Retry(func() error {
		db, err = sql.Open("mysql", dataSource)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// database migration
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("Could not instantiate driver: %s", err)
	}
	mig, err := migrate.NewWithDatabaseInstance(schemaFiles, "mysql", driver)
	if err != nil {
		log.Fatalf("Could not migrate database: %s", err)
	}
	mig.Up()

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestCreateAuthor(t *testing.T) {
	tests := []struct {
		scenario string
		input    sqlc.CreateAuthorParams
		expected sqlc.Author
	}{
		{
			scenario: "create user",
			input: sqlc.CreateAuthorParams{
				Name: "Brian Kernighan",
				Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
			},
			expected: sqlc.Author{
				Name: "Brian Kernighan",
				Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.scenario, func(t *testing.T) {
			queries := sqlc.New(db)

			// crete author
			ctx := context.Background()
			result, err := queries.CreateAuthor(ctx, tt.input)
			if err != nil {
				t.Error(err)
			}

			insertedAuthorID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// get author
			author, err := queries.GetAuthor(ctx, insertedAuthorID)
			if err != nil {
				t.Error(err)
			}

			if author.Name != tt.expected.Name {
				t.Fatalf("got=%v, want=%v", author.Name, tt.expected.Name)
			}
			if author.Bio != tt.expected.Bio {
				t.Fatalf("got=%v, want=%v", author.Bio, tt.expected.Bio)
			}
		})
	}
}
