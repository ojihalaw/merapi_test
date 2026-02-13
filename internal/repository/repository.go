package repository

import (
	"mertani_test/internal/utils"
	"strings"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	DB          *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}

func (r *Repository[T]) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("id = ?", id).Count(&total).Error
	return total, err
}

func (r *Repository[T]) FindById(db *gorm.DB, entity *T, id any) (*T, error) {
	if err := db.Where("id = ?", id).Take(&entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}


func (r *Repository[T]) FindAll(db *gorm.DB, entities *[]T, pagination *utils.PaginationRequest) (int64, error) {
	var total int64

	query := db.Model(new(T))

	if s, ok := any(new(T)).(utils.Searchable); ok && pagination.Search != "" {
		fields := s.SearchFields()
		conditions := make([]string, len(fields))
		args := make([]interface{}, len(fields))
		for i, f := range fields {
			conditions[i] = f + " ILIKE ?"
			args[i] = "%" + pagination.Search + "%"
		}
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	if pagination.OrderBy != "" {
		order := pagination.OrderBy
		if pagination.SortBy != "" {
			order += " " + pagination.SortBy
		}
		query = query.Order(order)
	}

	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}

	offset := (pagination.Page - 1) * pagination.Limit
	if err := query.Offset(offset).Limit(pagination.Limit).Find(entities).Error; err != nil {
		return 0, err
	}

	return total, nil
}
