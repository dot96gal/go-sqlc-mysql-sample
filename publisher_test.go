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

func TestCreatePublisher(t *testing.T) {
	publisherUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createPublisherParams sqlc.CreatePublisherParams
		}
		expected sqlc.Publisher
	}{
		{
			scenario: "create publisher",
			input: struct {
				createPublisherParams sqlc.CreatePublisherParams
			}{
				createPublisherParams: sqlc.CreatePublisherParams{
					Uuid: publisherUuid,
					Name: "publisher001",
				},
			},
			expected: sqlc.Publisher{
				Uuid: publisherUuid,
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
			err = queries.CreatePublisher(ctx, tt.input.createPublisherParams)
			if err != nil {
				t.Error(err)
			}

			// get publisher
			publisher, err := queries.GetPublisher(ctx, tt.input.createPublisherParams.Uuid)
			if err != nil {
				t.Error(err)
			}

			if publisher != tt.expected {
				t.Errorf("got=%v, want=%v", publisher, tt.expected)
			}
		})
	}
}

func TestUpdatePublisher(t *testing.T) {
	publisherUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createPublisherParams sqlc.CreatePublisherParams
			updatePublisherParams sqlc.UpdatePublisherParams
		}
		expected sqlc.Publisher
	}{
		{
			scenario: "update publisher",
			input: struct {
				createPublisherParams sqlc.CreatePublisherParams
				updatePublisherParams sqlc.UpdatePublisherParams
			}{
				createPublisherParams: sqlc.CreatePublisherParams{
					Uuid: publisherUuid,
					Name: "publisher001",
				},
				updatePublisherParams: sqlc.UpdatePublisherParams{
					Uuid: publisherUuid,
					Name: "Updated: publisher001",
				},
			},
			expected: sqlc.Publisher{
				Uuid: publisherUuid,
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
			err = queries.CreatePublisher(ctx, tt.input.createPublisherParams)
			if err != nil {
				t.Error(err)
			}

			// update publisher
			err = queries.UpdatePublisher(ctx, tt.input.updatePublisherParams)
			if err != nil {
				t.Error(err)
			}

			// get publisher
			publisher, err := queries.GetPublisher(ctx, tt.input.updatePublisherParams.Uuid)
			if err != nil {
				t.Error(err)
			}

			if publisher != tt.expected {
				t.Errorf("got=%v, want=%v", publisher, tt.expected)
			}
		})
	}
}

func TestDeletePublisher(t *testing.T) {
	publisherUuid := uuid.New()

	tests := []struct {
		scenario string
		input    struct {
			createPublisherParams sqlc.CreatePublisherParams
			deletePublisherUuid   uuid.UUID
		}
		expected error
	}{
		{
			scenario: "delete publisher",
			input: struct {
				createPublisherParams sqlc.CreatePublisherParams
				deletePublisherUuid   uuid.UUID
			}{
				createPublisherParams: sqlc.CreatePublisherParams{
					Uuid: publisherUuid,
					Name: "publisher001",
				},
				deletePublisherUuid: publisherUuid,
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

			// crete publisher
			ctx := context.Background()
			err = queries.CreatePublisher(ctx, tt.input.createPublisherParams)
			if err != nil {
				t.Error(err)
			}

			// delete publisher
			err = queries.DeletePublisher(ctx, tt.input.deletePublisherUuid)
			if err != nil {
				t.Error(err)
			}

			// get publisher
			_, err = queries.GetPublisher(ctx, tt.input.deletePublisherUuid)
			if err != tt.expected {
				t.Errorf("got=%v, want=%v", err, tt.expected)
			}
		})
	}
}

func TestListPublishers(t *testing.T) {
	publisherUuids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}

	tests := []struct {
		scenario string
		input    struct {
			createPublisherParamsList []sqlc.CreatePublisherParams
		}
		expected []sqlc.Publisher
	}{
		{
			scenario: "list publishers",
			input: struct{ createPublisherParamsList []sqlc.CreatePublisherParams }{
				createPublisherParamsList: []sqlc.CreatePublisherParams{
					{
						Uuid: publisherUuids[0],
						Name: "publisher001",
					},
					{
						Uuid: publisherUuids[1],
						Name: "publisher002",
					},
				},
			},
			expected: []sqlc.Publisher{
				{
					Uuid: publisherUuids[0],
					Name: "publisher001",
				},
				{
					Uuid: publisherUuids[1],
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
			for _, params := range tt.input.createPublisherParamsList {
				err := queries.CreatePublisher(ctx, params)
				if err != nil {
					t.Error(err)
				}
			}

			// list publishers
			publishers, err := queries.ListPublishers(ctx)
			if err != nil {
				t.Error(err)
			}

			sort.Slice(
				tt.expected,
				func(i, j int) bool {
					return tt.expected[i].Uuid.String() < tt.expected[j].Uuid.String()
				},
			)

			for i := range publishers {
				if publishers[i] != tt.expected[i] {
					t.Errorf("got=%v, want=%v", publishers[i].Name, tt.expected[i])
				}
			}
		})
	}
}

func TestGetPublisherBooks(t *testing.T) {
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
			getPublisherBooksUuid uuid.UUID
		}
		expected []sqlc.GetPublisherBooksRow
	}{
		{
			scenario: "get publisher books",
			input: struct {
				createPublisherParams sqlc.CreatePublisherParams
				createBookParamsList  []sqlc.CreateBookParams
				getPublisherBooksUuid uuid.UUID
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
				getPublisherBooksUuid: publisherUuid,
			},
			expected: []sqlc.GetPublisherBooksRow{
				{
					PublisherUuid: publisherUuid,
					PublisherName: "publisher001",
					BookUuid:      bookUuids[0],
					BookTitle:     "book001",
				},
				{
					PublisherUuid: publisherUuid,
					PublisherName: "publisher001",
					BookUuid:      bookUuids[1],
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
			err = queries.CreatePublisher(ctx, tt.input.createPublisherParams)
			if err != nil {
				t.Error(err)
			}

			// create book
			for _, params := range tt.input.createBookParamsList {
				err := queries.CreateBook(ctx, params)
				if err != nil {
					t.Error(err)
				}
			}

			// get publisher_books
			publisherBooks, err := queries.GetPublisherBooks(ctx, tt.input.getPublisherBooksUuid)
			if err != nil {
				t.Error(err)
			}

			sort.Slice(
				tt.expected,
				func(i, j int) bool {
					if tt.expected[i].PublisherUuid.String() < tt.expected[j].PublisherUuid.String() {
						return true
					} else if tt.expected[i].PublisherUuid.String() == tt.expected[j].PublisherUuid.String() {
						if tt.expected[i].BookUuid.String() < tt.expected[j].BookUuid.String() {
							return true
						}
					}
					return false
				},
			)

			for i := range publisherBooks {
				if publisherBooks[i] != tt.expected[i] {
					t.Errorf("got=%v, want=%v", publisherBooks, tt.expected)
				}
			}
		})
	}
}
