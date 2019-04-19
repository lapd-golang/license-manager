package helpers

import (
	"encoding/base64"
)

func GenerateLicenses(refs []string) []string {
	licenses := []string{}
	for _, license := range refs {
		licenses = append(licenses, base64.StdEncoding.EncodeToString([]byte(license)))
	}
	return licenses
}
