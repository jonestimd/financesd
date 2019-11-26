package companydao

import (
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/model"
)

func GetAll(db *gorm.DB) ([]*model.Company, error) {
	companies := make([]*model.Company, 0)
	err := db.Order("name").Find(&companies).Error
	return companies, err
}

func FindById(db *gorm.DB, id int64) (*model.Company, error) {
	var company model.Company
	err := db.First(&company, "id = ?", id).Error
	return &company, err
}

func FindByName(db *gorm.DB, name string) (*model.Company, error) {
	var company model.Company
	err := db.First(&company, "lower(name) = lower(?)", name).Error
	return &company, err
}
