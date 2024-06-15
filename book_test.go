package main

import (
	"context"
	"database/sql"
	"sort"
	"testing"

	"github.com/dot96gal/go-sqlc-mysql-sample/internal/sqlc"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
)

func TestCreateBook(t *testing.T) {
	publisherUuid := uuid.New()
	bookUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createPublisherPrams sqlc.CreatePublisherParams
			createBookParams     sqlc.CreateBookParams
		}
		expected sqlc.Book
	}{
		{
			scenario: "create book",
			input: struct {
				createPublisherPrams sqlc.CreatePublisherParams
				createBookParams     sqlc.CreateBookParams
			}{
				createPublisherPrams: sqlc.CreatePublisherParams{
					Uuid: publisherUuid,
					Name: "publisher001",
				},
				createBookParams: sqlc.CreateBookParams{
					Uuid:          bookUuid,
					Title:         "book001",
					PublisherUuid: publisherUuid,
				},
			},
			expected: sqlc.Book{
				Uuid:          bookUuid,
				Title:         "book001",
				PublisherUuid: publisherUuid,
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
			err = queries.CreatePublisher(ctx, tt.input.createPublisherPrams)
			if err != nil {
				t.Error(err)
			}

			// create book
			err = queries.CreateBook(ctx, tt.input.createBookParams)
			if err != nil {
				t.Error(err)
			}

			// get book
			book, err := queries.GetBook(ctx, tt.input.createBookParams.Uuid)
			if err != nil {
				t.Error(err)
			}

			if book != tt.expected {
				t.Errorf("got=%v, want=%v", book, tt.expected)
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	publisherUuid := uuid.New()
	bookUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createPublisherParams sqlc.CreatePublisherParams
			createBookParams      sqlc.CreateBookParams
			updateBookParams      sqlc.UpdateBookParams
		}
		expected sqlc.Book
	}{
		{
			scenario: "update book",
			input: struct {
				createPublisherParams sqlc.CreatePublisherParams
				createBookParams      sqlc.CreateBookParams
				updateBookParams      sqlc.UpdateBookParams
			}{
				createPublisherParams: sqlc.CreatePublisherParams{
					Uuid: publisherUuid,
					Name: "publisher001",
				},
				createBookParams: sqlc.CreateBookParams{
					Uuid:          bookUuid,
					Title:         "book001",
					PublisherUuid: publisherUuid,
				},
				updateBookParams: sqlc.UpdateBookParams{
					Uuid:  bookUuid,
					Title: "Updated: book001",
				},
			},
			expected: sqlc.Book{
				Uuid:          bookUuid,
				Title:         "Updated: book001",
				PublisherUuid: publisherUuid,
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
			err = queries.CreatePublisher(ctx, tt.input.createPublisherParams)
			if err != nil {
				t.Error(err)
			}

			// crete book
			err = queries.CreateBook(ctx, tt.input.createBookParams)
			if err != nil {
				t.Error(err)
			}

			// update book
			err = queries.UpdateBook(ctx, tt.input.updateBookParams)
			if err != nil {
				t.Error(err)
			}

			// get book
			book, err := queries.GetBook(ctx, tt.input.updateBookParams.Uuid)
			if err != nil {
				t.Error(err)
			}

			if book != tt.expected {
				t.Errorf("got=%v, want=%v", book, tt.expected)
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	publisherUuid := uuid.New()
	bookUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createPublisherParams sqlc.CreatePublisherParams
			createBookParams      sqlc.CreateBookParams
			deleteBookUuid        uuid.UUID
		}
		expected error
	}{
		{
			scenario: "delete book",
			input: struct {
				createPublisherParams sqlc.CreatePublisherParams
				createBookParams      sqlc.CreateBookParams
				deleteBookUuid        uuid.UUID
			}{
				createPublisherParams: sqlc.CreatePublisherParams{
					Uuid: publisherUuid,
					Name: "publisher001",
				},
				createBookParams: sqlc.CreateBookParams{
					Uuid:          bookUuid,
					Title:         "book001",
					PublisherUuid: publisherUuid,
				},
				deleteBookUuid: bookUuid,
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
			err = queries.CreatePublisher(ctx, tt.input.createPublisherParams)
			if err != nil {
				t.Error(err)
			}

			// crete book
			err = queries.CreateBook(ctx, tt.input.createBookParams)
			if err != nil {
				t.Error(err)
			}

			// delete book
			err = queries.DeleteBook(ctx, tt.input.deleteBookUuid)
			if err != nil {
				t.Error(err)
			}

			// get book
			_, err = queries.GetBook(ctx, tt.input.deleteBookUuid)
			if err != tt.expected {
				t.Errorf("got=%v, want=%v", err, tt.expected)
			}
		})
	}
}

func TestListBooks(t *testing.T) {
	publisherUuid := uuid.New()
	bookUuids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}

	tests := []struct {
		scenario string
		input    struct {
			createPublisherParams sqlc.CreatePublisherParams
			createBookParamsList  []sqlc.CreateBookParams
		}
		expected []sqlc.Book
	}{
		{
			scenario: "list books",
			input: struct {
				createPublisherParams sqlc.CreatePublisherParams
				createBookParamsList  []sqlc.CreateBookParams
			}{
				createPublisherParams: sqlc.CreatePublisherParams{
					Uuid: publisherUuid,
					Name: "publisher001",
				},
				createBookParamsList: []sqlc.CreateBookParams{
					{
						Uuid:          bookUuids[0],
						Title:         "book001",
						PublisherUuid: publisherUuid,
					},
					{
						Uuid:          bookUuids[1],
						Title:         "book002",
						PublisherUuid: publisherUuid,
					},
				},
			},
			expected: []sqlc.Book{
				{
					Uuid:          bookUuids[0],
					Title:         "book001",
					PublisherUuid: publisherUuid,
				},
				{
					Uuid:          bookUuids[1],
					Title:         "book002",
					PublisherUuid: publisherUuid,
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
			err = queries.CreatePublisher(ctx, tt.input.createPublisherParams)
			if err != nil {
				t.Error(err)
			}

			// crete book
			for _, params := range tt.input.createBookParamsList {
				err := queries.CreateBook(ctx, params)
				if err != nil {
					t.Error(err)
				}
			}

			// list books
			books, err := queries.ListBooks(ctx)
			if err != nil {
				t.Error(err)
			}

			sort.Slice(
				tt.expected,
				func(i, j int) bool {
					return tt.expected[i].Uuid.String() < tt.expected[j].Uuid.String()
				},
			)

			for i := range books {
				if books[i] != tt.expected[i] {
					t.Errorf("got=%v, want=%v", books[i], tt.expected[i])
				}
			}
		})
	}
}

func TestGetBookPublisher(t *testing.T) {
	publisherUuid := uuid.New()
	bookUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createPublisherParams sqlc.CreatePublisherParams
			createBookParams      sqlc.CreateBookParams
		}
		expected sqlc.GetBookPublisherRow
	}{
		{
			scenario: "get book publisher",
			input: struct {
				createPublisherParams sqlc.CreatePublisherParams
				createBookParams      sqlc.CreateBookParams
			}{
				createPublisherParams: sqlc.CreatePublisherParams{
					Uuid: publisherUuid,
					Name: "publisher001",
				},
				createBookParams: sqlc.CreateBookParams{
					Uuid:          bookUuid,
					Title:         "book001",
					PublisherUuid: publisherUuid,
				},
			},
			expected: sqlc.GetBookPublisherRow{
				BookUuid:      bookUuid,
				BookTitle:     "book001",
				PublisherUuid: publisherUuid,
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
			err = queries.CreatePublisher(ctx, tt.input.createPublisherParams)
			if err != nil {
				t.Error(err)
			}

			// create book
			err = queries.CreateBook(ctx, tt.input.createBookParams)
			if err != nil {
				t.Error(err)
			}

			// get book_publisher
			bookPublisher, err := queries.GetBookPublisher(ctx, tt.input.createBookParams.Uuid)
			if err != nil {
				t.Error(err)
			}

			if bookPublisher != tt.expected {
				t.Errorf("got=%v, want=%v", bookPublisher, tt.expected)
			}
		})
	}
}
