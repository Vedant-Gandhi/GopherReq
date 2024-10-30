package cookie

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"
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

var (
	ErrInvalidCookieFormat   = errors.New("Cookie format is invalid")
	ErrInvalidName           = errors.New("invalid cookie name")
	ErrInvalidValue          = errors.New("invalid cookie value")
	ErrInvalidDomain         = errors.New("invalid cookie domain")
	ErrInvalidPath           = errors.New("invalid cookie path")
	ErrInvalidExpires        = errors.New("invalid cookie expiration")
	ErrInvalidMaxAge         = errors.New("invalid cookie max-age")
	ErrInvalidSameSite       = errors.New("invalid cookie same-site attribute")
	ErrSecureRequiredForNone = errors.New("secure flag required when SameSite=None")
)

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

// isValidName checks if the cookie name follows RFC 6265 specs
func isValidName(name string) bool {
	if name == "" {
		return false
	}

	return strings.IndexFunc(name, func(r rune) bool {
		// Cookie names must not contain separator characters
		return unicode.IsSpace(r) || strings.ContainsRune("()<>@,;:\\\"/[]?={}", r)
	}) < 0
}

// isValidValue checks if the cookie value follows RFC 6265 specs
func isValidValue(value string) bool {
	if value == "" {
		return true // Empty values are allowed
	}

	return strings.IndexFunc(value, func(r rune) bool {
		// Cookie values must not contain separator characters or whitespace
		return r <= ' ' || r > '~' || strings.ContainsRune("(),/\\?@:;\"=", r)
	}) < 0
}

// isValidDomain checks if the cookie domain follows RFC 6265 specs
func isValidDomain(domain string) bool {
	if domain == "" {
		return true // Empty domain is allowed (defaults to current domain)
	}

	// Remove leading dot as per RFC 6265
	if domain[0] == '.' {
		domain = domain[1:]
	}

	// Basic domain name validation
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_.]+[a-zA-Z0-9]$`)
	return domainRegex.MatchString(domain) && !strings.Contains(domain, "..")
}

// isValidPath checks if the cookie path follows RFC 6265 specs
func isValidPath(path string) bool {
	if path == "" {
		return true // Empty path is allowed (defaults to current path)
	}

	// Path must start with "/"
	if !strings.HasPrefix(path, "/") {
		return false
	}

	// Check if path contains invalid characters
	return strings.IndexFunc(path, func(r rune) bool {
		return r <= ' ' || r > '~' || r == ';'
	}) < 0
}

func (c *Cookie) Validate() error {

	if !isValidName(c.Name) {
		return fmt.Errorf("%w: %s", ErrInvalidName, c.Name)
	}

	if !isValidValue(c.Value) {
		return fmt.Errorf("%w: %s", ErrInvalidValue, c.Value)
	}

	if !isValidDomain(c.Domain) {
		return fmt.Errorf("%w: %s", ErrInvalidDomain, c.Domain)
	}

	if !isValidPath(c.Path) {
		return fmt.Errorf("%w: %s", ErrInvalidPath, c.Path)
	}

	switch c.SameSite {
	case SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode:

	default:
		return fmt.Errorf("%w: %d", ErrInvalidSameSite, c.SameSite)
	}

	if c.SameSite == SameSiteNoneMode && !c.Secure {
		return ErrSecureRequiredForNone
	}

	if strings.Contains(c.Value, "http://") || strings.Contains(c.Value, "https://") {
		_, err := url.Parse(c.Value)
		if err != nil {
			return fmt.Errorf("%w: invalid URL in cookie value", ErrInvalidValue)
		}
	}

	return nil
}
func ParseRequestCookie(cookie string) (c Cookie, err error) {

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

}

func (l *CookieList) Exists(key string) (exists bool) {

	_, exists = l.Get(key)

	return
}
