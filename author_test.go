package main

import (
	"context"
	"database/sql"
	"sort"
	"testing"

	"github.com/dot96gal/go-sqlc-sample/internal/sqlc"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
)

func TestCreateAuthor(t *testing.T) {
	authorUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createAuthorParams sqlc.CreateAuthorParams
		}
		expected sqlc.Author
	}{
		{
			scenario: "create author",
			input: struct {
				createAuthorParams sqlc.CreateAuthorParams
			}{
				createAuthorParams: sqlc.CreateAuthorParams{
					Uuid: authorUuid,
					Name: "author001",
					Bio:  sql.NullString{String: "author001", Valid: true},
				},
			},
			expected: sqlc.Author{
				Uuid: authorUuid,
				Name: "author001",
				Bio:  sql.NullString{String: "author001", Valid: true},
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

			// get author
			author, err := queries.GetAuthor(ctx, tt.input.createAuthorParams.Uuid)
			if err != nil {
				t.Error(err)
			}

			if author != tt.expected {
				t.Errorf("got=%v, want=%v", author, tt.expected)
			}
		})
	}
}

func TestUpdateAuthor(t *testing.T) {
	authorUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createAuthorParams sqlc.CreateAuthorParams
			updateAuthorParams sqlc.UpdateAuthorParams
		}
		expected sqlc.Author
	}{
		{
			scenario: "update author",
			input: struct {
				createAuthorParams sqlc.CreateAuthorParams
				updateAuthorParams sqlc.UpdateAuthorParams
			}{
				createAuthorParams: sqlc.CreateAuthorParams{
					Uuid: authorUuid,
					Name: "author001",
					Bio:  sql.NullString{String: "author001", Valid: true},
				},
				updateAuthorParams: sqlc.UpdateAuthorParams{
					Name: "Updated: author001",
					Bio:  sql.NullString{String: "Updated: author001", Valid: true},
					Uuid: authorUuid,
				},
			},
			expected: sqlc.Author{
				Uuid: authorUuid,
				Name: "Updated: author001",
				Bio:  sql.NullString{String: "Updated: author001", Valid: true},
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

			// update author
			err = queries.UpdateAuthor(ctx, tt.input.updateAuthorParams)
			if err != nil {
				t.Error(err)
			}

			// get author
			author, err := queries.GetAuthor(ctx, tt.input.updateAuthorParams.Uuid)
			if err != nil {
				t.Error(err)
			}

			if author != tt.expected {
				t.Errorf("got=%v, want=%v", author, tt.expected)
			}
		})
	}
}

func TestDeleteAuthor(t *testing.T) {
	authorUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createAuthorParams sqlc.CreateAuthorParams
			deleteAuthorUuid   uuid.UUID
		}
		expected error
	}{
		{
			scenario: "delete author",
			input: struct {
				createAuthorParams sqlc.CreateAuthorParams
				deleteAuthorUuid   uuid.UUID
			}{
				createAuthorParams: sqlc.CreateAuthorParams{
					Uuid: authorUuid,
					Name: "author001",
					Bio:  sql.NullString{String: "author001", Valid: true},
				},
				deleteAuthorUuid: authorUuid,
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

			// delete author
			err = queries.DeleteAuthor(ctx, tt.input.deleteAuthorUuid)
			if err != nil {
				t.Error(err)
			}

			// get author
			_, err = queries.GetAuthor(ctx, tt.input.deleteAuthorUuid)
			if err != tt.expected {
				t.Errorf("got=%v, want=%v", err, tt.expected)
			}
		})
	}
}

func TestListAuthors(t *testing.T) {
	authorUuids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}

	tests := []struct {
		scenario string
		input    struct {
			createAuthorParamsList []sqlc.CreateAuthorParams
		}
		expected []sqlc.Author
	}{
		{
			scenario: "list authors",
			input: struct{ createAuthorParamsList []sqlc.CreateAuthorParams }{
				createAuthorParamsList: []sqlc.CreateAuthorParams{
					{
						Uuid: authorUuids[0],
						Name: "author001",
						Bio:  sql.NullString{String: "author001", Valid: true},
					},
					{
						Uuid: authorUuids[1],
						Name: "author002",
						Bio:  sql.NullString{String: "author002", Valid: true},
					},
				},
			},
			expected: []sqlc.Author{
				{
					Uuid: authorUuids[0],
					Name: "author001",
					Bio:  sql.NullString{String: "author001", Valid: true},
				},
				{
					Uuid: authorUuids[1],
					Name: "author002",
					Bio:  sql.NullString{String: "author002", Valid: true},
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

			// list authors
			authors, err := queries.ListAuthors(ctx)
			if err != nil {
				t.Error(err)
			}

			sort.Slice(
				tt.expected,
				func(i, j int) bool {
					return tt.expected[i].Uuid.String() < tt.expected[j].Uuid.String()
				},
			)

			for i := range authors {
				if authors[i] != tt.expected[i] {
					t.Errorf("got=%v, want=%v", authors[i], tt.expected[i])
				}
			}
		})
	}
}
