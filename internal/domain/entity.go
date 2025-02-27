package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ParamEntityID string = "entity-id"
	ParamSearch   string = "search"
)

type Entity struct {
	ID        uuid.UUID  `json:"id" param:"id" query:"id" header:"id" validate:"required" gorm:"primarykey"`
	CreatedAt time.Time  `json:"created_at" db:"created_at" gorm:"->"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at" validate:"required" gorm:"->"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

type IEntity interface {
	GetUpdatedAt() time.Time
	GetID() uuid.UUID
}

func (e *Entity) GetUpdatedAt() time.Time {
	return e.UpdatedAt.UTC()
}

func (e *Entity) GetID() uuid.UUID {
	return e.ID
}

func (e *Entity) IsSame(other IEntity) bool {
	return e.ID == other.GetID()
}

func (e *Entity) IsSameVersion(other IEntity) bool {
	return e.ID == other.GetID() && e.GetUpdatedAt().Equal(other.GetUpdatedAt())
}

func (e *Entity) BeforeCreate(tx *gorm.DB) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	e.ID = uuid
	return nil
}
