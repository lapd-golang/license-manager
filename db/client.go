package db

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/sevren/test/models"
	log "github.com/sirupsen/logrus"
)

type Dao struct {
	DB *gorm.DB
}

//StoreUsedLicenses - Stores each license recieved from the rabbitMQ exchange
func (d *Dao) StoreUsedLicenses(license string) error {

	l := models.Used_licenses{License: license}
	var res *gorm.DB
	db := d.DB.Table("used_licenses")

	// Check if the database already contains that license
	if db.Where("license = ?", license).RecordNotFound() {
		res = db.Create(l)

		if res.Error != nil {
			log.Error(res.Error)
			return res.Error
		}
	}

	return nil
}

//GetLicenses - Retrieves the Users licenses from the database
func (d *Dao) GetLicenses(user string) ([]string, error) {
	u := models.User_licenses{}
	if err := d.DB.Where(&models.User_licenses{Username: user}).First(&u).Error; err != nil {
		return nil, errors.New("User: " + user + "has no licenses")
	}
	return u.Lics, nil
}
