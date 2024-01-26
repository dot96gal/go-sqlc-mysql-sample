package main

import (
	"context"
	"database/sql"
	"sort"
	"testing"

	"github.com/dot96gal/go-sqlc-sample/internal/sqlc"
	"github.com/google/uuid"
)

func TestCreateAuthorBook(t *testing.T) {
	authorUuid := uuid.New()
	publisherUuid := uuid.New()
	bookUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createAuthorParams     sqlc.CreateAuthorParams
			createPublisherParams  sqlc.CreatePublisherParams
			createBookParams       sqlc.CreateBookParams
			createAuthorBookParams sqlc.CreateAuthorBookParams
		}
		expected sqlc.AuthorBook
	}{
		{
			scenario: "create author_book",
			input: struct {
				createAuthorParams     sqlc.CreateAuthorParams
				createPublisherParams  sqlc.CreatePublisherParams
				createBookParams       sqlc.CreateBookParams
				createAuthorBookParams sqlc.CreateAuthorBookParams
			}{
				createAuthorParams: sqlc.CreateAuthorParams{
					Uuid: authorUuid,
					Name: "author001",
					Bio:  sql.NullString{String: "author001", Valid: true},
				},
				createPublisherParams: sqlc.CreatePublisherParams{
					Uuid: publisherUuid,
					Name: "publisher001",
				},
				createBookParams: sqlc.CreateBookParams{
					Uuid:          bookUuid,
					Title:         "book001",
					PublisherUuid: publisherUuid,
				},
				createAuthorBookParams: sqlc.CreateAuthorBookParams{
					AuthorUuid: authorUuid,
					BookUuid:   bookUuid,
				},
			},
			expected: sqlc.AuthorBook{
				AuthorUuid: authorUuid,
				BookUuid:   bookUuid,
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
			err = queries.CreateAuthor(ctx, tt.input.createAuthorParams)
			if err != nil {
				t.Error(err)
			}

			// create publisher
			err = queries.CreatePublisher(ctx, tt.input.createPublisherParams)
			if err != nil {
				t.Error(err)
			}

			// crete book
			err = queries.CreateBook(ctx, tt.input.createBookParams)
			if err != nil {
				t.Error(err)
			}

			// create author_book
			err = queries.CreateAuthorBook(ctx, tt.input.createAuthorBookParams)
			if err != nil {
				t.Error(err)
			}

			// get author_book
			authorBook, err := queries.GetAuthorBook(
				ctx,
				sqlc.GetAuthorBookParams{
					AuthorUuid: tt.input.createAuthorBookParams.AuthorUuid,
					BookUuid:   tt.input.createAuthorBookParams.BookUuid,
				},
			)
			if err != nil {
				t.Error(err)
			}

			if authorBook != tt.expected {
				t.Errorf("got=%v, want=%v", authorBook, tt.expected)
			}
		})
	}
}

func TestDeleteAuthorBook(t *testing.T) {
	authorUuid := uuid.New()
	publisherUuid := uuid.New()
	bookUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createAuthorParams     sqlc.CreateAuthorParams
			createPublisherParams  sqlc.CreatePublisherParams
			createBookParams       sqlc.CreateBookParams
			createAuthorBookParams sqlc.CreateAuthorBookParams
			deleteAuthorBookParams sqlc.DeleteAuthorBookParams
		}
		expected error
	}{
		{
			scenario: "delete author_book",
			input: struct {
				createAuthorParams     sqlc.CreateAuthorParams
				createPublisherParams  sqlc.CreatePublisherParams
				createBookParams       sqlc.CreateBookParams
				createAuthorBookParams sqlc.CreateAuthorBookParams
				deleteAuthorBookParams sqlc.DeleteAuthorBookParams
			}{
				createAuthorParams: sqlc.CreateAuthorParams{
					Uuid: authorUuid,
					Name: "author001",
					Bio:  sql.NullString{String: "author001", Valid: true},
				},
				createPublisherParams: sqlc.CreatePublisherParams{
					Uuid: publisherUuid,
					Name: "publisher001",
				},
				createBookParams: sqlc.CreateBookParams{
					Uuid:          bookUuid,
					Title:         "book001",
					PublisherUuid: publisherUuid,
				},
				createAuthorBookParams: sqlc.CreateAuthorBookParams{
					AuthorUuid: authorUuid,
					BookUuid:   bookUuid,
				},
				deleteAuthorBookParams: sqlc.DeleteAuthorBookParams{
					AuthorUuid: authorUuid,
					BookUuid:   bookUuid,
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
			err = queries.CreateAuthor(ctx, tt.input.createAuthorParams)
			if err != nil {
				t.Error(err)
			}

			// create publisher
			err = queries.CreatePublisher(ctx, tt.input.createPublisherParams)
			if err != nil {
				t.Error(err)
			}

			// crete book
			err = queries.CreateBook(ctx, tt.input.createBookParams)
			if err != nil {
				t.Error(err)
			}

			// create author_book
			err = queries.CreateAuthorBook(ctx, tt.input.createAuthorBookParams)
			if err != nil {
				t.Error(err)
			}

			// delete author_book
			err = queries.DeleteAuthorBook(ctx, tt.input.deleteAuthorBookParams)
			if err != nil {
				t.Error(err)
			}

			// get author_book
			_, err = queries.GetAuthorBook(
				ctx,
				sqlc.GetAuthorBookParams{
					AuthorUuid: tt.input.deleteAuthorBookParams.AuthorUuid,
					BookUuid:   tt.input.deleteAuthorBookParams.BookUuid,
				},
			)
			if err != tt.expected {
				t.Errorf("got=%v, want=%v", err, tt.expected)
			}
		})
	}
}

func TestListAuthorBooks(t *testing.T) {
	authorUuids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}
	publisherUuid := uuid.New()
	bookUuids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}

	tests := []struct {
		scenario string
		input    struct {
			createAuthorParamsList     []sqlc.CreateAuthorParams
			createPublisherParams      sqlc.CreatePublisherParams
			createBookParamsList       []sqlc.CreateBookParams
			createAuthorBookParamsList []sqlc.CreateAuthorBookParams
		}
		expected []sqlc.ListAuthorBooksRow
	}{
		{
			scenario: "list author_books",
			input: struct {
				createAuthorParamsList     []sqlc.CreateAuthorParams
				createPublisherParams      sqlc.CreatePublisherParams
				createBookParamsList       []sqlc.CreateBookParams
				createAuthorBookParamsList []sqlc.CreateAuthorBookParams
			}{
				createAuthorParamsList: []sqlc.CreateAuthorParams{
					{
						Uuid: authorUuids[0],
						Name: "author001",
						Bio:  sql.NullString{String: "author001", Valid: true},
					},
					{
						Uuid: authorUuids[1],
						Name: "author002",
						Bio:  sql.NullString{String: "author001", Valid: true},
					},
				},
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
				createAuthorBookParamsList: []sqlc.CreateAuthorBookParams{
					{
						AuthorUuid: authorUuids[0],
						BookUuid:   bookUuids[0],
					},
					{
						AuthorUuid: authorUuids[0],
						BookUuid:   bookUuids[1],
					},
					{
						AuthorUuid: authorUuids[1],
						BookUuid:   bookUuids[0],
					},
					{
						AuthorUuid: authorUuids[1],
						BookUuid:   bookUuids[1],
					},
				},
			},
			expected: []sqlc.ListAuthorBooksRow{
				{
					AuthorUuid: authorUuids[0],
					AuthorName: "author001",
					AuthorBio:  sql.NullString{String: "author001", Valid: true},
					BookUuid:   bookUuids[0],
					BookTitle:  "book001",
				},
				{
					AuthorUuid: authorUuids[0],
					AuthorName: "author001",
					AuthorBio:  sql.NullString{String: "author001", Valid: true},
					BookUuid:   bookUuids[1],
					BookTitle:  "book002",
				},
				{
					AuthorUuid: authorUuids[1],
					AuthorName: "author002",
					AuthorBio:  sql.NullString{String: "author001", Valid: true},
					BookUuid:   bookUuids[0],
					BookTitle:  "book001",
				},
				{
					AuthorUuid: authorUuids[1],
					AuthorName: "author002",
					AuthorBio:  sql.NullString{String: "author001", Valid: true},
					BookUuid:   bookUuids[1],
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
			for _, params := range tt.input.createAuthorParamsList {
				err := queries.CreateAuthor(ctx, params)
				if err != nil {
					t.Error(err)
				}
			}

			// create publisher
			err = queries.CreatePublisher(ctx, tt.input.createPublisherParams)
			if err != nil {
				t.Error(err)
			}

			// crete book
			for _, params := range tt.input.createBookParamsList {
				err = queries.CreateBook(ctx, params)
				if err != nil {
					t.Error(err)
				}
			}

			// create author_book
			for _, params := range tt.input.createAuthorBookParamsList {
				err = queries.CreateAuthorBook(ctx, params)
				if err != nil {
					t.Error(err)
				}
			}

			// list author_books
			authorBooks, err := queries.ListAuthorBooks(ctx)
			if err != nil {
				t.Error(err)
			}

			sort.Slice(
				tt.expected,
				func(i, j int) bool {
					if tt.expected[i].AuthorUuid.String() < tt.expected[j].AuthorUuid.String() {
						return true
					} else if tt.expected[i].AuthorUuid.String() == tt.expected[j].AuthorUuid.String() {
						if tt.expected[i].BookUuid.String() < tt.expected[j].BookUuid.String() {
							return true
						}
					}
					return false
				},
			)

			for i := range authorBooks {
				if authorBooks[i] != tt.expected[i] {
					t.Errorf("got=%v, want=%v", authorBooks[i], tt.expected[i])
				}
			}
		})
	}
}
