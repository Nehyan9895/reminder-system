package repository

import "gorm.io/gorm"

type GormRepo struct {
	DB *gorm.DB
}

func NewGormRepo(db *gorm.DB) *GormRepo {
	return &GormRepo{DB: db}
}
