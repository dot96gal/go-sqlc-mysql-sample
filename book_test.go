package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dot96gal/go-sqlc-sample/internal/sqlc"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func TestCreateBook(t *testing.T) {
	tests := []struct {
		scenario string
		input    sqlc.CreateBookParams
		expected sqlc.Book
	}{
		{
			scenario: "create book",
			input: sqlc.CreateBookParams{
				Title: "book001",
			},
			expected: sqlc.Book{
				Title: "book001",
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

			// crete publisher
			ctx := context.Background()
			publisher, err := queries.CreatePublisher(ctx, "publisher001")
			if err != nil {
				t.Error(err)
			}

			insertedPublisherID, err := publisher.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// create book
			tt.input.PublisherID = insertedPublisherID
			result, err := queries.CreateBook(ctx, tt.input)
			if err != nil {
				t.Error(err)
			}

			insertedBookID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// get book
			book, err := queries.GetBook(ctx, insertedBookID)
			if err != nil {
				t.Error(err)
			}

			if book.Title != tt.expected.Title {
				t.Errorf("got=%v, want=%v", book.Title, tt.expected.Title)
			}
			if book.PublisherID != insertedBookID {
				t.Errorf("got=%v, want=%v", book.PublisherID, insertedBookID)
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	tests := []struct {
		scenario string
		input    sqlc.UpdateBookParams
		expected sqlc.Book
	}{
		{
			scenario: "update book",
			input: sqlc.UpdateBookParams{
				Title: "Updated: book001",
			},
			expected: sqlc.Book{
				Title: "Updated: book001",
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

			// create publisher
			ctx := context.Background()
			publisher, err := queries.CreatePublisher(ctx, "publisher001")
			if err != nil {
				t.Error(err)
			}

			insertedPublisherID, err := publisher.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// crete book
			input := sqlc.CreateBookParams{Title: "", PublisherID: insertedPublisherID}
			result, err := queries.CreateBook(ctx, input)
			if err != nil {
				t.Error(err)
			}

			insertedBookID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// update book
			tt.input.ID = insertedBookID
			err = queries.UpdateBook(ctx, tt.input)
			if err != nil {
				t.Error(err)
			}

			// get book
			book, err := queries.GetBook(ctx, insertedBookID)
			if err != nil {
				t.Error(err)
			}

			if book.Title != tt.expected.Title {
				t.Errorf("got=%v, want=%v", book.Title, tt.expected.Title)
			}
			if book.PublisherID != insertedBookID {
				t.Errorf("got=%v, want=%v", book.PublisherID, insertedBookID)
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	tests := []struct {
		scenario string
		input    sqlc.CreateBookParams
		expected error
	}{
		{
			scenario: "delete book",
			input: sqlc.CreateBookParams{
				Title: "book001",
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

			// create publisher
			ctx := context.Background()
			publisher, err := queries.CreatePublisher(ctx, "publisher001")
			if err != nil {
				t.Error(err)
			}

			insertedPublisherID, err := publisher.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// crete book
			tt.input.PublisherID = insertedPublisherID
			result, err := queries.CreateBook(ctx, tt.input)
			if err != nil {
				t.Error(err)
			}

			insertedBookID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// delete book
			err = queries.DeleteBook(ctx, insertedBookID)
			if err != nil {
				t.Error(err)
			}

			// get book
			_, err = queries.GetBook(ctx, insertedBookID)
			if err != tt.expected {
				t.Errorf("got=%v, want=%v", err, tt.expected)
			}
		})
	}
}

func TestListBook(t *testing.T) {
	tests := []struct {
		scenario string
		input    []sqlc.CreateBookParams
		expected []sqlc.Book
	}{
		{
			scenario: "list book",
			input: []sqlc.CreateBookParams{
				{
					Title: "book001",
				},
				{
					Title: "book002",
				},
			},
			expected: []sqlc.Book{
				{
					Title: "book001",
				},
				{
					Title: "book002",
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

			// create publisher
			ctx := context.Background()
			publisher, err := queries.CreatePublisher(ctx, "publisher001")
			if err != nil {
				t.Error(err)
			}

			insertedPublisherID, err := publisher.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// crete book
			for _, input := range tt.input {
				input.PublisherID = insertedPublisherID
				_, err := queries.CreateBook(ctx, input)
				if err != nil {
					t.Error(err)
				}
			}

			// list book
			results, err := queries.ListBooks(ctx)
			if err != nil {
				t.Error(err)
			}

			for i := 0; i < len(tt.expected); i++ {
				if results[i].Title != tt.expected[i].Title {
					t.Errorf("got=%v, want=%v", results[i].Title, tt.expected[i].Title)
				}
				if results[i].PublisherID != insertedPublisherID {
					t.Errorf("got=%v, want=%v", results[i].PublisherID, insertedPublisherID)
				}
			}
		})
	}
}

func TestGetBookPublisher(t *testing.T) {
	tests := []struct {
		scenario string
		input    struct {
			bookParams      sqlc.CreateBookParams
			publisherParams string
		}
		expected sqlc.GetBookPublisherRow
	}{
		{
			scenario: "get book publisher",
			input: struct {
				bookParams      sqlc.CreateBookParams
				publisherParams string
			}{
				bookParams: sqlc.CreateBookParams{
					Title: "book001",
				},
				publisherParams: "publisher001",
			},
			expected: sqlc.GetBookPublisherRow{
				BookTitle:     "book001",
				PublisherName: "publisher001",
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

			// crete publisher
			ctx := context.Background()
			publisher, err := queries.CreatePublisher(ctx, tt.input.publisherParams)
			if err != nil {
				t.Error(err)
			}

			insertedPublisherID, err := publisher.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// create book
			tt.input.bookParams.PublisherID = insertedPublisherID
			book, err := queries.CreateBook(ctx, tt.input.bookParams)
			if err != nil {
				t.Error(err)
			}

			insertedBookID, err := book.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// get book publisher
			bookPublisher, err := queries.GetBookPublisher(ctx, insertedBookID)
			if err != nil {
				t.Error(err)
			}

			tt.expected.BookID = insertedBookID
			tt.expected.PublisherID = insertedPublisherID

			if bookPublisher != tt.expected {
				t.Errorf("got=%v, want=%v", bookPublisher, tt.expected)
			}
		})
	}
}
