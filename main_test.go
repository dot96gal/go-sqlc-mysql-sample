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
		log.Fatalf("Could not instantiate migrate: %s", err)
	}
	err = mig.Up()
	if err != nil {
		log.Fatalf("Could not migrate database: %s", err)
	}

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

			// test with transaction
			tx, err := db.Begin()
			if err != nil {
				t.Error(err)
			}
			t.Cleanup(func() {
				err = tx.Rollback()
				if err != nil {
					t.Error(err)
				}
			})

			queries = queries.WithTx(tx)

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
				t.Errorf("got=%v, want=%v", author.Name, tt.expected.Name)
			}
			if author.Bio != tt.expected.Bio {
				t.Errorf("got=%v, want=%v", author.Bio, tt.expected.Bio)
			}
		})
	}
}

func TestUpdateAuthor(t *testing.T) {
	tests := []struct {
		scenario string
		input    sqlc.UpdateAuthorParams
		expected sqlc.Author
	}{
		{
			scenario: "update user",
			input: sqlc.UpdateAuthorParams{
				Name: "Updated: Brian Kernighan",
				Bio:  sql.NullString{String: "Updated: Co-author of The C Programming Language and The Go Programming Language", Valid: true},
			},
			expected: sqlc.Author{
				Name: "Updated: Brian Kernighan",
				Bio:  sql.NullString{String: "Updated: Co-author of The C Programming Language and The Go Programming Language", Valid: true},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.scenario, func(t *testing.T) {
			queries := sqlc.New(db)

			// test with transaction
			tx, err := db.Begin()
			if err != nil {
				t.Error(err)
			}
			t.Cleanup(func() {
				err = tx.Rollback()
				if err != nil {
					t.Error(err)
				}
			})

			queries = queries.WithTx(tx)

			// crete author
			ctx := context.Background()
			input := sqlc.CreateAuthorParams{Name: "", Bio: sql.NullString{String: "", Valid: true}}
			result, err := queries.CreateAuthor(ctx, input)
			if err != nil {
				t.Error(err)
			}

			insertedAuthorID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// update author
			tt.input.ID = insertedAuthorID
			err = queries.UpdateAuthor(ctx, tt.input)
			if err != nil {
				t.Error(err)
			}

			// get author
			author, err := queries.GetAuthor(ctx, insertedAuthorID)
			if err != nil {
				t.Error(err)
			}

			if author.Name != tt.expected.Name {
				t.Errorf("got=%v, want=%v", author.Name, tt.expected.Name)
			}
			if author.Bio != tt.expected.Bio {
				t.Errorf("got=%v, want=%v", author.Bio, tt.expected.Bio)
			}
		})
	}
}

func TestDeleteAuthor(t *testing.T) {
	tests := []struct {
		scenario string
		input    sqlc.CreateAuthorParams
		expected error
	}{
		{
			scenario: "delete user",
			input: sqlc.CreateAuthorParams{
				Name: "Brian Kernighan",
				Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
			},
			expected: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.scenario, func(t *testing.T) {
			queries := sqlc.New(db)

			// test with transaction
			tx, err := db.Begin()
			if err != nil {
				t.Error(err)
			}
			t.Cleanup(func() {
				err = tx.Rollback()
				if err != nil {
					t.Error(err)
				}
			})

			queries = queries.WithTx(tx)

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

			// delete author
			err = queries.DeleteAuthor(ctx, insertedAuthorID)
			if err != nil {
				t.Error(err)
			}

			// get author
			_, err = queries.GetAuthor(ctx, insertedAuthorID)
			if err != tt.expected {
				t.Errorf("got=%v, want=%v", err, tt.expected)
			}
		})
	}
}

func TestListAuthor(t *testing.T) {
	tests := []struct {
		scenario string
		input    []sqlc.CreateAuthorParams
		expected []sqlc.Author
	}{
		{
			scenario: "list user",
			input: []sqlc.CreateAuthorParams{
				{
					Name: "hoge",
					Bio:  sql.NullString{String: "hoge", Valid: true},
				},
				{
					Name: "fuga",
					Bio:  sql.NullString{String: "fuga", Valid: true},
				},
			},
			expected: []sqlc.Author{
				{
					Name: "hoge",
					Bio:  sql.NullString{String: "hoge", Valid: true},
				},
				{
					Name: "fuga",
					Bio:  sql.NullString{String: "fuga", Valid: true},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.scenario, func(t *testing.T) {
			queries := sqlc.New(db)

			// test with transaction
			tx, err := db.Begin()
			if err != nil {
				t.Error(err)
			}
			t.Cleanup(func() {
				err = tx.Rollback()
				if err != nil {
					t.Error(err)
				}
			})

			queries = queries.WithTx(tx)

			// crete author
			ctx := context.Background()
			for _, input := range tt.input {
				_, err := queries.CreateAuthor(ctx, input)
				if err != nil {
					t.Error(err)
				}
			}

			// list author
			results, err := queries.ListAuthors(ctx)
			if err != nil {
				t.Error(err)
			}

			for i := 0; i < len(tt.expected); i++ {
				if results[i].Name != tt.expected[i].Name {
					t.Errorf("got=%v, want=%v", results[i].Name, tt.expected[i].Name)
				}
				if results[i].Bio != tt.expected[i].Bio {
					t.Errorf("got=%v, want=%v", results[i].Bio, tt.expected[i].Bio)
				}
			}
		})
	}
}
