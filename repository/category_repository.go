package repository

import (
	"errors"
	"gorm.io/gorm"
	"hannibal/gin-try/common"
	"hannibal/gin-try/model"
	"log"
)

//重复代码抽离
type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository() CategoryRepository {
	return CategoryRepository{
		DB: common.GetDB(),
	}
}
func (c CategoryRepository) Create(name string) (*model.Category, error) {
	category := model.Category{
		Name: name,
	}
	if result := c.DB.Create(&category); result.Error != nil || result.RowsAffected == 0 {
		//errors.Unwrap(fmt.Errorf("创建失败: %w", result.Error))
		if result.Error != nil {
			panic(result.Error)
		} else {
			return nil, errors.New("创建失败")
		}
	} else {
		return &category, nil
	}
}
func (c CategoryRepository) Update(category model.Category, name string) (*model.Category, error) {
	if result := c.DB.Model(&category).Update("name", name); result.Error != nil || result.RowsAffected == 0 {
		if result.Error != nil {
			panic(result.Error)
		} else {
			return nil, errors.New("修改失败")
		}
	} else {
		log.Printf("%+v", result)
		return &category, nil
	}
}

func (c CategoryRepository) SelectById(id int) (*model.Category, error) {
	var category model.Category
	if result := c.DB.First(&category, id); result.Error != nil {
		//errors.Unwrap(fmt.Errorf("查询失败: %w", result.Error))
		if result.Error != nil {
			panic(result.Error)
		} else {
			return nil, errors.New("查询失败")
		}
	} else {
		return &category, nil
	}
}

func (c CategoryRepository) DeleteById(id int) error {
	if result := c.DB.Delete(model.Category{}, id); result.Error != nil || result.RowsAffected == 0 {
		//errors.Unwrap(fmt.Errorf("删除失败: %w", result.Error))
		if result.Error != nil {
			panic(result.Error)
		} else {
			return errors.New("删除失败")
		}
	} else {
		return nil
	}
}
