package db

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

const reconn = 1
const table = "user_licenses"

// Connect - Attempts to connect to the database
func Connect(db string) (*Dao, error) {

	log.WithFields(log.Fields{
		"database": db,
	}).Infof("Connecting to database")

	ticker := time.NewTicker(time.Duration(5) * time.Second)
	for i := 0; i < reconn; i++ {
		<-ticker.C
		dbc, err := gorm.Open("sqlite3", db)
		if err != nil {
			log.Error(err)
			continue
		}

		dao := &Dao{
			DB: dbc.Table(table),
		}

		return dao, nil
	}

	return nil, errors.New("Unable to connect to the sqlite database")

}
