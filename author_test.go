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
			scenario: "update author",
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
			scenario: "delete author",
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
			scenario: "list author",
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
