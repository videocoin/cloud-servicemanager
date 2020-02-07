package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Base contains the base model definition.
type Base struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate updates the timestamp for `CreatedAt`, `UpdatedAt`.
func (b *Base) BeforeCreate(scope *gorm.Scope) error {
	t := time.Now()
	scope.Set("CreatedAt", t)
	scope.Set("UpdatedAt", t)
	return nil
}

// BeforeUpdate updates the timestamp for `UpdatedAt`.
func (b *Base) BeforeUpdate(scope *gorm.Scope) error {
	scope.Set("UpdatedAt", time.Now())
	return nil
}
