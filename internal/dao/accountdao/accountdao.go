package accountdao

import (
	"github.com/jinzhu/gorm"
	"github.com/jonestimd/financesd/internal/model"
)

// GetAccounts gets all of the accounts.
func GetAll(db *gorm.DB) ([]*model.Account, error) {
	accounts := make([]*model.Account, 0)
	err := db.Debug().Order("name").Find(&accounts).Error
	return accounts, err
}
