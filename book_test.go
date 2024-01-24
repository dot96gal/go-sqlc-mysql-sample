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
		input    struct {
			publisherName    string
			createBookParams sqlc.CreateBookParams
		}
		expected sqlc.Book
	}{
		{
			scenario: "create book",
			input: struct {
				publisherName    string
				createBookParams sqlc.CreateBookParams
			}{
				publisherName: "publisher001",
				createBookParams: sqlc.CreateBookParams{
					Title: "book001",
				},
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
			publisher, err := queries.CreatePublisher(ctx, tt.input.publisherName)
			if err != nil {
				t.Error(err)
			}

			publisherID, err := publisher.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// create book
			tt.input.createBookParams.PublisherID = publisherID
			result, err := queries.CreateBook(ctx, tt.input.createBookParams)
			if err != nil {
				t.Error(err)
			}

			bookID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// get book
			book, err := queries.GetBook(ctx, bookID)
			if err != nil {
				t.Error(err)
			}

			tt.expected.ID = bookID
			tt.expected.PublisherID = publisherID
			if book != tt.expected {
				t.Errorf("got=%v, want=%v", book, tt.expected)
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	tests := []struct {
		scenario string
		input    struct {
			publisherName    string
			createBookParams sqlc.CreateBookParams
			updateBookParams sqlc.UpdateBookParams
		}
		expected sqlc.Book
	}{
		{
			scenario: "update book",
			input: struct {
				publisherName    string
				createBookParams sqlc.CreateBookParams
				updateBookParams sqlc.UpdateBookParams
			}{
				publisherName: "publisher001",
				createBookParams: sqlc.CreateBookParams{
					Title: "book001",
				},
				updateBookParams: sqlc.UpdateBookParams{
					Title: "Updated: book001",
				},
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
			result, err := queries.CreatePublisher(ctx, tt.input.publisherName)
			if err != nil {
				t.Error(err)
			}

			publisherID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// crete book
			tt.input.createBookParams.PublisherID = publisherID
			result, err = queries.CreateBook(ctx, tt.input.createBookParams)
			if err != nil {
				t.Error(err)
			}

			bookID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// update book
			tt.input.updateBookParams.ID = bookID
			err = queries.UpdateBook(ctx, tt.input.updateBookParams)
			if err != nil {
				t.Error(err)
			}

			// get book
			book, err := queries.GetBook(ctx, bookID)
			if err != nil {
				t.Error(err)
			}

			tt.expected.ID = bookID
			tt.expected.PublisherID = publisherID
			if book != tt.expected {
				t.Errorf("got=%v, want=%v", book, tt.expected)
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	tests := []struct {
		scenario string
		input    struct {
			publisherName    string
			createBookParams sqlc.CreateBookParams
		}
		expected error
	}{
		{
			scenario: "delete book",
			input: struct {
				publisherName    string
				createBookParams sqlc.CreateBookParams
			}{
				publisherName: "publisher001",
				createBookParams: sqlc.CreateBookParams{
					Title: "book001",
				}},
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
			result, err := queries.CreatePublisher(ctx, tt.input.publisherName)
			if err != nil {
				t.Error(err)
			}

			publisherID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// crete book
			tt.input.createBookParams.PublisherID = publisherID
			result, err = queries.CreateBook(ctx, tt.input.createBookParams)
			if err != nil {
				t.Error(err)
			}

			bookID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// delete book
			err = queries.DeleteBook(ctx, bookID)
			if err != nil {
				t.Error(err)
			}

			// get book
			_, err = queries.GetBook(ctx, bookID)
			if err != tt.expected {
				t.Errorf("got=%v, want=%v", err, tt.expected)
			}
		})
	}
}

func TestListBooks(t *testing.T) {
	tests := []struct {
		scenario string
		input    struct {
			publisherName    string
			createBookParams []sqlc.CreateBookParams
		}
		expected []sqlc.Book
	}{
		{
			scenario: "list books",
			input: struct {
				publisherName    string
				createBookParams []sqlc.CreateBookParams
			}{
				publisherName: "publisher001",
				createBookParams: []sqlc.CreateBookParams{
					{
						Title: "book001",
					},
					{
						Title: "book002",
					},
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
			result, err := queries.CreatePublisher(ctx, tt.input.publisherName)
			if err != nil {
				t.Error(err)
			}

			publisherID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// crete book
			bookIDs := []int64{}
			for _, params := range tt.input.createBookParams {
				params.PublisherID = publisherID
				result, err := queries.CreateBook(ctx, params)
				if err != nil {
					t.Error(err)
				}

				bookID, err := result.LastInsertId()
				if err != nil {
					t.Error(err)
				}

				bookIDs = append(bookIDs, bookID)
			}

			// list books
			books, err := queries.ListBooks(ctx)
			if err != nil {
				t.Error(err)
			}

			for i := range bookIDs {
				tt.expected[i].ID = bookIDs[i]
				tt.expected[i].PublisherID = publisherID
				if books[i] != tt.expected[i] {
					t.Errorf("got=%v, want=%v", books[i], tt.expected[i])
				}
			}
		})
	}
}

func TestGetBookPublisher(t *testing.T) {
	tests := []struct {
		scenario string
		input    struct {
			createBookParams sqlc.CreateBookParams
			publisherName    string
		}
		expected sqlc.GetBookPublisherRow
	}{
		{
			scenario: "get book publisher",
			input: struct {
				createBookParams sqlc.CreateBookParams
				publisherName    string
			}{
				createBookParams: sqlc.CreateBookParams{
					Title: "book001",
				},
				publisherName: "publisher001",
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
			result, err := queries.CreatePublisher(ctx, tt.input.publisherName)
			if err != nil {
				t.Error(err)
			}

			publisherID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// create book
			tt.input.createBookParams.PublisherID = publisherID
			result, err = queries.CreateBook(ctx, tt.input.createBookParams)
			if err != nil {
				t.Error(err)
			}

			bookID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// get book publisher
			bookPublisher, err := queries.GetBookPublisher(ctx, bookID)
			if err != nil {
				t.Error(err)
			}

			tt.expected.BookID = bookID
			tt.expected.PublisherID = publisherID

			if bookPublisher != tt.expected {
				t.Errorf("got=%v, want=%v", bookPublisher, tt.expected)
			}
		})
	}
}
