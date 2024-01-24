package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dot96gal/go-sqlc-sample/internal/sqlc"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func TestCreateAuthor(t *testing.T) {
	tests := []struct {
		scenario string
		input    sqlc.CreateAuthorParams
		expected sqlc.Author
	}{
		{
			scenario: "create author",
			input: sqlc.CreateAuthorParams{
				Name: "author001",
				Bio:  sql.NullString{String: "author001", Valid: true},
			},
			expected: sqlc.Author{
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
			result, err := queries.CreateAuthor(ctx, tt.input)
			if err != nil {
				t.Error(err)
			}

			authorID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// get author
			author, err := queries.GetAuthor(ctx, authorID)
			if err != nil {
				t.Error(err)
			}

			tt.expected.ID = authorID
			if author != tt.expected {
				t.Errorf("got=%v, want=%v", author, tt.expected)
			}
		})
	}
}

func TestUpdateAuthor(t *testing.T) {
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
					Name: "author001",
					Bio:  sql.NullString{String: "author001", Valid: true},
				},
				updateAuthorParams: sqlc.UpdateAuthorParams{
					Name: "Updated: author001",
					Bio:  sql.NullString{String: "Updated: author001", Valid: true},
				},
			},
			expected: sqlc.Author{
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
			result, err := queries.CreateAuthor(ctx, tt.input.createAuthorParams)
			if err != nil {
				t.Error(err)
			}

			authorID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// update author
			tt.input.updateAuthorParams.ID = authorID
			err = queries.UpdateAuthor(ctx, tt.input.updateAuthorParams)
			if err != nil {
				t.Error(err)
			}

			// get author
			author, err := queries.GetAuthor(ctx, authorID)
			if err != nil {
				t.Error(err)
			}

			tt.expected.ID = authorID
			if author != tt.expected {
				t.Errorf("got=%v, want=%v", author, tt.expected)
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
			scenario: "delete author",
			input: sqlc.CreateAuthorParams{
				Name: "author001",
				Bio:  sql.NullString{String: "author001", Valid: true},
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

			authorID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// delete author
			err = queries.DeleteAuthor(ctx, authorID)
			if err != nil {
				t.Error(err)
			}

			// get author
			_, err = queries.GetAuthor(ctx, authorID)
			if err != tt.expected {
				t.Errorf("got=%v, want=%v", err, tt.expected)
			}
		})
	}
}

func TestListAuthors(t *testing.T) {
	tests := []struct {
		scenario string
		input    []sqlc.CreateAuthorParams
		expected []sqlc.Author
	}{
		{
			scenario: "list authors",
			input: []sqlc.CreateAuthorParams{
				{
					Name: "author001",
					Bio:  sql.NullString{String: "author001", Valid: true},
				},
				{
					Name: "author002",
					Bio:  sql.NullString{String: "author002", Valid: true},
				},
			},
			expected: []sqlc.Author{
				{
					Name: "author001",
					Bio:  sql.NullString{String: "author001", Valid: true},
				},
				{
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
			authorIDs := []int64{}
			for _, input := range tt.input {
				result, err := queries.CreateAuthor(ctx, input)
				if err != nil {
					t.Error(err)
				}

				authorID, err := result.LastInsertId()
				if err != nil {
					t.Error(err)
				}

				authorIDs = append(authorIDs, authorID)
			}

			// list authors
			authors, err := queries.ListAuthors(ctx)
			if err != nil {
				t.Error(err)
			}

			for i := range authorIDs {
				tt.expected[i].ID = authorIDs[i]
				if authors[i] != tt.expected[i] {
					t.Errorf("got=%v, want=%v", authors[i], tt.expected[i])
				}
			}
		})
	}
}
