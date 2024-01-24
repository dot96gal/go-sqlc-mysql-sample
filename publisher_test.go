package main

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dot96gal/go-sqlc-sample/internal/sqlc"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func TestCreatePublisher(t *testing.T) {
	tests := []struct {
		scenario string
		input    string
		expected sqlc.Publisher
	}{
		{
			scenario: "create publisher",
			input:    "publisher001",
			expected: sqlc.Publisher{
				Name: "publisher001",
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
			result, err := queries.CreatePublisher(ctx, tt.input)
			if err != nil {
				t.Error(err)
			}

			publisherID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// get publisher
			publisher, err := queries.GetPublisher(ctx, publisherID)
			if err != nil {
				t.Error(err)
			}

			tt.expected.ID = publisherID
			if publisher != tt.expected {
				t.Errorf("got=%v, want=%v", publisher, tt.expected)
			}
		})
	}
}

func TestUpdatePublisher(t *testing.T) {
	tests := []struct {
		scenario string
		input    struct {
			publisherName         string
			updatePublisherParams sqlc.UpdatePublisherParams
		}
		expected sqlc.Publisher
	}{
		{
			scenario: "update publisher",
			input: struct {
				publisherName         string
				updatePublisherParams sqlc.UpdatePublisherParams
			}{
				publisherName: "publisher001",
				updatePublisherParams: sqlc.UpdatePublisherParams{
					Name: "Updated: publisher001",
				},
			},
			expected: sqlc.Publisher{
				Name: "Updated: publisher001",
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

			// update publisher
			tt.input.updatePublisherParams.ID = publisherID
			err = queries.UpdatePublisher(ctx, tt.input.updatePublisherParams)
			if err != nil {
				t.Error(err)
			}

			// get publisher
			publisher, err := queries.GetPublisher(ctx, publisherID)
			if err != nil {
				t.Error(err)
			}

			tt.expected.ID = publisherID
			if publisher != tt.expected {
				t.Errorf("got=%v, want=%v", publisher, tt.expected)
			}
		})
	}
}

func TestDeletePublisher(t *testing.T) {
	tests := []struct {
		scenario string
		input    string
		expected error
	}{
		{
			scenario: "delete publisher",
			input:    "publisher001",
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

			// crete publisher
			ctx := context.Background()
			result, err := queries.CreatePublisher(ctx, tt.input)
			if err != nil {
				t.Error(err)
			}

			publisherID, err := result.LastInsertId()
			if err != nil {
				t.Error(err)
			}

			// delete publisher
			err = queries.DeletePublisher(ctx, publisherID)
			if err != nil {
				t.Error(err)
			}

			// get publisher
			_, err = queries.GetPublisher(ctx, publisherID)
			if err != tt.expected {
				t.Errorf("got=%v, want=%v", err, tt.expected)
			}
		})
	}
}

func TestListPublishers(t *testing.T) {
	tests := []struct {
		scenario string
		input    []string
		expected []sqlc.Publisher
	}{
		{
			scenario: "list publishers",
			input: []string{
				"publisher001",
				"publisher002",
			},
			expected: []sqlc.Publisher{
				{
					Name: "publisher001",
				},
				{
					Name: "publisher002",
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

			// crete publisher
			ctx := context.Background()
			publisherIDs := []int64{}
			for _, input := range tt.input {
				result, err := queries.CreatePublisher(ctx, input)
				if err != nil {
					t.Error(err)
				}

				publisherID, err := result.LastInsertId()
				if err != nil {
					t.Error(err)
				}

				publisherIDs = append(publisherIDs, publisherID)
			}

			// list publishers
			publishers, err := queries.ListPublishers(ctx)
			if err != nil {
				t.Error(err)
			}

			for i := range publisherIDs {
				tt.expected[i].ID = publisherIDs[i]
				if publishers[i] != tt.expected[i] {
					t.Errorf("got=%v, want=%v", publishers[i].Name, tt.expected[i])
				}
			}
		})
	}
}

func TestGetPublisherBooks(t *testing.T) {
	tests := []struct {
		scenario string
		input    struct {
			publisherName    string
			createBookParams []sqlc.CreateBookParams
		}
		expected []sqlc.GetPublisherBooksRow
	}{
		{
			scenario: "get publisher books",
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
			expected: []sqlc.GetPublisherBooksRow{
				{
					PublisherName: "publisher001",
					BookTitle:     "book001",
				},
				{
					PublisherName: "publisher001",
					BookTitle:     "book002",
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
			bookIDs := []int64{}
			for _, params := range tt.input.createBookParams {
				params.PublisherID = publisherID
				book, err := queries.CreateBook(ctx, params)
				if err != nil {
					t.Error(err)
				}

				bookID, err := book.LastInsertId()
				if err != nil {
					t.Error(err)
				}

				bookIDs = append(bookIDs, bookID)
			}

			// get publisher books
			publisherBooks, err := queries.GetPublisherBooks(ctx, publisherID)
			if err != nil {
				t.Error(err)
			}

			for i := range bookIDs {
				tt.expected[i].PublisherID = publisherID
				tt.expected[i].BookID = bookIDs[i]
				if publisherBooks[i] != tt.expected[i] {
					t.Errorf("got=%v, want=%v", publisherBooks, tt.expected)
				}
			}
		})
	}
}
