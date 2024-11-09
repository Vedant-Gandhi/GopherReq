package common

import "strings"

type StatusCode int

const CRLF = "\r\n"

type HttpMethod string

const (
	Get    HttpMethod = "GET"
	Post   HttpMethod = "POST"
	Put    HttpMethod = "PUT"
	Delete HttpMethod = "DELETE"
)

// This function converts the key name to canonical name.
func GetCanonicalName(key string) (canonical string) {
	key = strings.Trim(key, " ")
	splitStr := strings.Split(key, "-")

	// Convert the each key to upper case.
	for index, stringPart := range splitStr {
		stringPart = strings.ToUpper(string(stringPart[0])) + stringPart[1:]
		splitStr[index] = stringPart
	}

	return strings.Join(splitStr, "-")

}
