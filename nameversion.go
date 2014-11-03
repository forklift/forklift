package main

import (
	"errors"
	"regexp"
)

var (
	nameversionRegexp = regexp.MustCompile(`^([a-zA-Z0-9]+[-a-zA-Z0-9]*[a-zA-Z0-9]+)` +
		`(?:-(` +
		`(?:[0-9]+(?:\.[0-9]*){0,2})?` +
		`(?:\-(?:pre|alpha|beta|rc)?(?:.[0-9]+)?)?` +
		`))?` +
		`(?:\+([0-9A-Za-z.]+))?$`)
)

type NameVersion struct {
	Location string
	Name     string
	Version  string
	Meta     string
}

func NewNameVersion(raw string) (NameVersion, error) {
	nv := NameVersion{}

	parts := nameversionRegexp.FindStringSubmatch(raw)

	if len(parts) < 1 {
		return nv, errors.New("Invalid Name-Version.")
	}

	if parts[2] == "" {
		parts[2] = "*"
	}

	nv.Name = parts[1]
	nv.Version = parts[2]
	nv.Meta = parts[3]
	return nv, nil
}
