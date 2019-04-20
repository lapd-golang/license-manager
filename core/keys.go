package core

import (
	"encoding/base64"
	"fmt"

	"github.com/speps/go-hashids"
)

// Generates Challenge1 licenses
func GenerateLicenses(refs []string) []string {
	licenses := []string{}
	for _, license := range refs {
		licenses = append(licenses, base64.StdEncoding.EncodeToString([]byte(license)))
	}
	return licenses
}

// Generates Challenge3 licenses
// Generates unique uuid for the user and the license text
func GenerateBetterLicenses(user string, refs []string) []string {

	licenses := []string{}
	hd := hashids.NewData()
	hd.MinLength = 30

	for _, license := range refs {
		hd.Salt = fmt.Sprintf("%s:%s", user, license)
		h, _ := hashids.NewWithData(hd)
		e, _ := h.Encode([]int{45, 434, 1313, 99})

		licenses = append(licenses, e)
	}
	return licenses
}
