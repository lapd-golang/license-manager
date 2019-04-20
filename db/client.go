package db

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/sevren/test/models"
)

type Dao struct {
	DB *gorm.DB
}

func (d *Dao) StoreUsedLicenses(license string) {

}

//GetLicenses - Retrieves the Users licenses from the database
func (d *Dao) GetLicenses(user string) ([]string, error) {
	u := models.User_licenses{}
	if err := d.DB.Where(&models.User_licenses{Username: user}).First(&u).Error; err != nil {
		return nil, errors.New("User: " + user + "has no licenses")
	}
	return u.Lics, nil
}
