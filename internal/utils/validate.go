package utils

import "regexp"

var addregexp = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

func IsAddress(v string) bool {
	return addregexp.MatchString(v)
}

var numberRx = regexp.MustCompile(`^[0-9]+$`)

func IsNumber(p string) bool {
	return numberRx.MatchString(p)
}

var hexRx = regexp.MustCompile(`^[0-9]+$`)

func IsHex(p string) bool {
	return hexRx.MatchString(p)
}
