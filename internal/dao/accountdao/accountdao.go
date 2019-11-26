package accountdao

import (
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/model"
)

// GetAccounts gets all of the accounts.
func GetAll(db *gorm.DB) ([]*model.Account, error) {
	accounts := make([]*model.Account, 0)
	err := db.Order("name").Find(&accounts).Error
	return accounts, err
}

func FindById(db *gorm.DB, id int) (*model.Account, error) {
	var account model.Account
	err := db.First(&account, "id = ?", id).Error
	return &account, err
}

func FindByName(db *gorm.DB, name string) (*model.Account, error) {
	var account model.Account
	err := db.First(&account, "lower(name) = lower(?)", name).Error
	return &account, err
}
