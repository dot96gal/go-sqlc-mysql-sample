package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/dot96gal/go-sqlc-mysql-sample/internal/sqlc"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func run() error {
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPass := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_TCP_PORT")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDatabase)

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return err
	}

	queries := sqlc.New(db)

	ctx := context.Background()
	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	authorUuid := uuid.New()

	err = queries.CreateAuthor(ctx, sqlc.CreateAuthorParams{
		Uuid: authorUuid,
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		return err
	}

	author, err := queries.GetAuthor(ctx, authorUuid)
	if err != nil {
		return err
	}
	log.Println(author)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
