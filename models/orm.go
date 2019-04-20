package models

import (
	"fmt"
	"strings"
)

// Gorm forces us to have the same struct name as the table name :/
type Used_licenses struct {
	License string `gorm:"type:text;column:license"`
}

type User_licenses struct {
	Username string   `gorm:"column:username"`
	Password string   `gorm:"column:password"`
	Lics     Licenses `gorm:"column:licenses"`
}

// Custom type defined for Licenses column
type Licenses []string

// Gorm - SQLlite3 apparently scans into []uint8 - need to convert
// This will be comma seperated values
func (v *Licenses) Scan(src interface{}) error {
	temp, ok := src.([]uint8)
	if !ok {
		return fmt.Errorf("Unable to convert %v of %T to CustomType", src, src)
	}
	*v = strings.Split(string(temp), ",")
	return nil
}
