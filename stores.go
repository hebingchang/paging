package paging

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// -----------------------------------------------------------------------------
// Interfaces
// -----------------------------------------------------------------------------

// Store is a store.
type Store interface {
	PaginateOffset(limit, offset int64, count *int64) error
	PaginateCursor(limit int64, cursor interface{}, fieldName string, reverse bool, hasnext *bool) error
	GetItems() interface{}
}

// -----------------------------------------------------------------------------
// GORM Store
// -----------------------------------------------------------------------------

// GORMStore is the store for GORM ORM.
type GORMStore struct {
	db    *gorm.DB
	items interface{}
}

// NewGORMStore returns a new GORM store instance.
func NewGORMStore(db *gorm.DB, items interface{}) (*GORMStore, error) {
	return &GORMStore{
		db:    db,
		items: items,
	}, nil
}

// GetItems return the current result
func (s *GORMStore) GetItems() interface{} {
	return s.items
}

// PaginateOffset paginates items from the store and update page instance.
func (s *GORMStore) PaginateOffset(limit, offset int64, count *int64) error {
	q := s.db
	q = q.Limit(int(limit))
	q = q.Offset(int(offset))


	if q = q.Find(s.items); q.Error != nil {
		return q.Error
	}

	q = q.Limit(-1)
	q = q.Offset(-1)

	if err := q.Count(count).Error; err != nil {
		return err
	}

	return nil
}

// PaginateCursor paginates items from the store and update page instance for cursor pagination system.
// cursor can be an ID or a date (time.Time)
func (s *GORMStore) PaginateCursor(limit int64, cursor interface{}, fieldName string, reverse bool, hasnext *bool) error {
	q := s.db

	q = q.Limit(limit + 1)

	if reverse {
		q = q.Where(fmt.Sprintf("%s < ?", fieldName), cursor)
	} else {
		q = q.Where(fmt.Sprintf("%s > ?", fieldName), cursor)
	}

	err := q.Find(s.items).Error
	if err != nil {
		return err
	}

	len := getLen(s.items)
	if int64(len) <= limit {
		*hasnext = false
		return nil
	}

	*hasnext = true
	_, s.items = popLastElement(s.items)
	return nil
}
