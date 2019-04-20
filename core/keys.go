package core

import (
	"encoding/base64"
)

// Generates Challenge1 licenses
func GenerateLicenses(refs []string) []string {
	licenses := []string{}
	for _, license := range refs {
		licenses = append(licenses, base64.StdEncoding.EncodeToString([]byte(license)))
	}
	return licenses
}
