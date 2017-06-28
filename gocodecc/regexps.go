package gocodecc

import (
	"regexp"
)

var mentionPeopleReg = regexp.MustCompile(`@[a-zA-Z0-9_]{5,20}?\s+`)
