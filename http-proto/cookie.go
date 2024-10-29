package httpproto

import (
	"errors"
	"strings"
	"time"
)

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
	SameSite   SameSite
	Raw        string
	Unparsed   []string // Raw text of unparsed attribute-value pairs
}

type CookieList struct {
	cookies map[string]Cookie
}

// SameSite represents the SameSite attribute of a cookie.
type SameSite int

const (
	SameSiteDefaultMode SameSite = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
	SameSiteNoneMode
)

var ErrInvalidCookieFormat = errors.New("Cookie format is not valid.")

// String returns the string representation of SameSite attribute
func (s SameSite) String() string {
	switch s {
	case SameSiteDefaultMode:
		return "Default"
	case SameSiteLaxMode:
		return "Lax"
	case SameSiteStrictMode:
		return "Strict"
	case SameSiteNoneMode:
		return "None"
	}
	return "Unknown"
}

func parseRequestCookie(cookie string) (c Cookie, err error) {

	splits := strings.SplitN(cookie, "=", 2)

	if len(splits) != 2 {
		err = ErrInvalidCookieFormat
		return
	}

	c = Cookie{
		Name:     splits[0],
		Value:    splits[1],
		Unparsed: splits,
	}

	return

}

func NewCookieList() CookieList {
	c := CookieList{
		cookies: make(map[string]Cookie),
	}

	return c
}

func (l *CookieList) Get(key string) (value Cookie, exists bool) {
	value = l.cookies[key]
	exists = len(value.Name) != 0

	return

}

func (l *CookieList) Add(c Cookie) {
	l.cookies[c.Name] = c
	return

}

func (l *CookieList) Exists(key string) (exists bool) {

	_, exists = l.Get(key)

	return
}
