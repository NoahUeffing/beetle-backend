package domain_test

import (
	"beetle/internal/domain"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestEntityIsSame_Success(t *testing.T) {
	is := is.New(t)
	e1 := &domain.Entity{
		ID: uuid.New(),
	}
	e2 := &domain.User{
		Entity: *e1,
	}

	is.True(e1.IsSame(e2)) // should be the same entity
}

func TestEntityIsSame_Failure(t *testing.T) {
	is := is.New(t)
	e1 := &domain.Entity{
		ID: uuid.New(),
	}
	e2 := &domain.User{
		Entity: domain.Entity{
			ID: uuid.New(),
		},
	}

	is.Equal(e1.IsSame(e2), false) // should be a different entity
}

func TestEntityIsSameVersion_Success(t *testing.T) {
	is := is.New(t)
	e1 := &domain.Entity{
		ID:        uuid.New(),
		UpdatedAt: time.Now(),
	}
	e2 := &domain.User{
		Entity: *e1,
	}

	is.True(e1.IsSameVersion(e2)) // should be the same version
}

func TestEntityIsSameVersion_Failure(t *testing.T) {
	is := is.New(t)
	e1 := &domain.Entity{
		ID:        uuid.New(),
		UpdatedAt: time.Now(),
	}
	e2 := &domain.User{
		Entity: domain.Entity{
			ID:        e1.ID,
			UpdatedAt: e1.UpdatedAt.Add(time.Hour),
		},
	}

	is.True(e1.IsSame(e2))                // should be the same entity
	is.Equal(e1.IsSameVersion(e2), false) // should be a different version
}
