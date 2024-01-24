package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dot96gal/go-sqlc-sample/internal/sqlc"
)

func TestCreateAuthorBook(t *testing.T) {
	tests := []struct {
		scenario string
		input    struct {
			createAuthorParams sqlc.CreateAuthorParams
			publisherName      string
			createBookParams   sqlc.CreateBookParams
		}
		expected sqlc.AuthorBook
	}{
		{
			scenario: "create author_book",
			input: struct {
				createAuthorParams sqlc.CreateAuthorParams
				publisherName      string
				createBookParams   sqlc.CreateBookParams
			}{
				createAuthorParams: sqlc.CreateAuthorParams{
					Name: "author001",
					Bio:  sql.NullString{String: "author001", Valid: true},
				},
				publisherName: "publisher001",
				createBookParams: sqlc.CreateBookParams{
					Title: "book001",
				},
			},
			expected: sqlc.AuthorBook{},
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
			result, err := queries.CreateAuthor(ctx, tt.input.createAuthorParams)
			if err != nil {
				t.Error(err)
			}

			authorID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// create publisher
			result, err = queries.CreatePublisher(ctx, tt.input.publisherName)
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

			// create author_book
			err = queries.CreateAuthorBook(ctx, sqlc.CreateAuthorBookParams{AuthorID: authorID, BookID: bookID})
			if err != nil {
				t.Error(err)
			}

			// get author_book
			authorBook, err := queries.GetAuthorBook(ctx, sqlc.GetAuthorBookParams{AuthorID: authorID, BookID: bookID})
			if err != nil {
				t.Error(err)
			}

			tt.expected.AuthorID = authorID
			tt.expected.BookID = bookID
			if authorBook != tt.expected {
				t.Errorf("got=%v, want=%v", authorBook, tt.expected)
			}
		})
	}
}

func TestDeleteAuthorBook(t *testing.T) {
	tests := []struct {
		scenario string
		input    struct {
			createAuthorParams sqlc.CreateAuthorParams
			publisherName      string
			createBookParams   sqlc.CreateBookParams
		}
		expected error
	}{
		{
			scenario: "delete author_book",
			input: struct {
				createAuthorParams sqlc.CreateAuthorParams
				publisherName      string
				createBookParams   sqlc.CreateBookParams
			}{
				createAuthorParams: sqlc.CreateAuthorParams{
					Name: "author001",
					Bio:  sql.NullString{String: "author001", Valid: true},
				},
				publisherName: "publisher001",
				createBookParams: sqlc.CreateBookParams{
					Title: "book001",
				},
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
			result, err := queries.CreateAuthor(ctx, tt.input.createAuthorParams)
			if err != nil {
				t.Error(err)
			}

			authorID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// create publisher
			result, err = queries.CreatePublisher(ctx, tt.input.publisherName)
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

			// create author_book
			err = queries.CreateAuthorBook(ctx, sqlc.CreateAuthorBookParams{AuthorID: authorID, BookID: bookID})
			if err != nil {
				t.Error(err)
			}

			// delete author_book
			err = queries.DeleteAuthorBook(ctx, sqlc.DeleteAuthorBookParams{AuthorID: authorID, BookID: bookID})
			if err != nil {
				t.Error(err)
			}

			// get author_book
			_, err = queries.GetAuthorBook(ctx, sqlc.GetAuthorBookParams{AuthorID: authorID, BookID: bookID})
			if err != tt.expected {
				t.Errorf("got=%v, want=%v", err, tt.expected)
			}
		})
	}
}

func TestListAuthorBooks(t *testing.T) {
	tests := []struct {
		scenario string
		input    struct {
			createAuthorParams []sqlc.CreateAuthorParams
			publisherName      string
			createBookParams   []sqlc.CreateBookParams
		}
		expected []sqlc.ListAuthorBooksRow
	}{
		{
			scenario: "list author_books",
			input: struct {
				createAuthorParams []sqlc.CreateAuthorParams
				publisherName      string
				createBookParams   []sqlc.CreateBookParams
			}{
				createAuthorParams: []sqlc.CreateAuthorParams{
					{
						Name: "author001",
						Bio:  sql.NullString{String: "author001", Valid: true},
					},
					{
						Name: "author002",
						Bio:  sql.NullString{String: "author001", Valid: true},
					},
				},
				publisherName: "publisher001",
				createBookParams: []sqlc.CreateBookParams{
					{Title: "book001"},
					{Title: "book002"},
				},
			},
			expected: []sqlc.ListAuthorBooksRow{
				{
					AuthorName: "author001",
					AuthorBio:  sql.NullString{String: "author001", Valid: true},
					BookTitle:  "book001",
				},
				{
					AuthorName: "author001",
					AuthorBio:  sql.NullString{String: "author001", Valid: true},
					BookTitle:  "book002",
				},
				{
					AuthorName: "author002",
					AuthorBio:  sql.NullString{String: "author001", Valid: true},
					BookTitle:  "book001",
				},
				{
					AuthorName: "author002",
					AuthorBio:  sql.NullString{String: "author001", Valid: true},
					BookTitle:  "book002",
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
			authorIDs := []int64{}
			for _, params := range tt.input.createAuthorParams {
				result, err := queries.CreateAuthor(ctx, params)
				if err != nil {
					t.Error(err)
				}

				authorID, err := result.LastInsertId()
				if err != nil {
					t.Error(err)
				}

				authorIDs = append(authorIDs, authorID)
			}

			// create publisher
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
				result, err = queries.CreateBook(ctx, params)
				if err != nil {
					t.Error(err)
				}

				bookID, err := result.LastInsertId()
				if err != nil {
					t.Error(err)
				}

				bookIDs = append(bookIDs, bookID)
			}

			// create author_book
			for _, authorID := range authorIDs {
				for _, bookID := range bookIDs {
					err = queries.CreateAuthorBook(ctx, sqlc.CreateAuthorBookParams{AuthorID: authorID, BookID: bookID})
					if err != nil {
						t.Error(err)
					}
				}
			}

			// list author_books
			authorBooks, err := queries.ListAuthorBooks(ctx)
			if err != nil {
				t.Error(err)
			}

			i := 0
			for _, authorID := range authorIDs {
				for _, bookID := range bookIDs {
					tt.expected[i].AuthorID = authorID
					tt.expected[i].BookID = bookID
					if authorBooks[i] != tt.expected[i] {
						t.Errorf("got=%v, want=%v", authorBooks[i], tt.expected[i])
					}
					i++
				}
			}
		})
	}
}
