package storage

import (
	"context"
	"testing"
	"time"

	"github.com/tanveerprottoy/backend-structure-go/internal/api/product"
	"github.com/tanveerprottoy/backend-structure-go/internal/api/product/postgres"
)

func TestProduct(t *testing.T) {
	// init storage
	s := postgres.NewStorage(db)
	// Mock data
	n := time.Now().Unix()

	description := "test description"

	var id string

	// test create
	t.Run("create", func(t *testing.T) {
		// t.Parallel()

		dto := &product.CreateDTO{
			Name:        "Test",
			Description: &description,
		}

		// call create
		lastID, err := s.Create(context.Background(), dto)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if lastID == "" {
			t.Errorf("expected id %s, got empty", lastID)
		}

		id = lastID
	})

	t.Run(("read many"), func(t *testing.T) {
		// t.Parallel()
		_, err := s.ReadMany(context.Background(), 10, 0)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})

	t.Run(("read one"), func(t *testing.T) {
		// t.Parallel()
		got, err := s.ReadOne(context.Background(), id)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if got.ID != id {
			t.Errorf("expected id %s, got %s", id, got.ID)
		}
	})

	t.Run(("update"), func(t *testing.T) {
		// t.Parallel()

		dto := &product.UpdateDTO{
			Name: "test 2",
		}

		_, err := s.Update(context.Background(), id, dto)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})

	t.Run(("delete"), func(t *testing.T) {
		_, err := s.Delete(context.Background(), id, n)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
}
