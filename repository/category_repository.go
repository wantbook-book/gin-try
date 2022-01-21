package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"hannibal/gin-try/model"
)

//重复代码抽离
type CategoryRepository struct {
	DB *gorm.DB
}

func (c CategoryRepository) Create(name string) (*model.Category, error) {
	category := model.Category{
		Name: name,
	}
	if result := c.DB.Create(&category); result.Error != nil || result.RowsAffected == 0 {

		return nil, errors.Unwrap(fmt.Errorf("创建失败: %w", result.Error))
	} else {
		return &category, nil
	}
}
