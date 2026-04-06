package surepository

import (
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("resource not found")

type SuperUtilRepository[T any] interface {
	Create(item *T) error
	GetByID(id uint) (*T, error)
	GetAll() ([]T, error)
	Update(item *T) error
	DeleteByID(id uint) error
}

type superUtilRepository[T any] struct {
	db *gorm.DB
}

func NewSuperUtilRepository[T any](db *gorm.DB) SuperUtilRepository[T] {
	return &superUtilRepository[T]{db: db}
}

func (r *superUtilRepository[T]) Create(item *T) error {
	return r.db.Create(item).Error
}

func (r *superUtilRepository[T]) GetByID(id uint) (*T, error) {
	var item T
	if err := r.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}
	return &item, nil
}

func (r *superUtilRepository[T]) GetAll() ([]T, error) {
	var items []T
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *superUtilRepository[T]) Update(item *T) error {
	if err := r.db.Save(item).Error; err != nil {
		return err
	}
	return nil
}

func (r *superUtilRepository[T]) DeleteByID(id uint) error {
	var item T
	if err := r.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}

		return err
	}
	if err := r.db.Delete(&item).Error; err != nil {
		return err
	}
	return nil
}
