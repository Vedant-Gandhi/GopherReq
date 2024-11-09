package httpproto

import (
	"http-v1_1/http-proto/common"
	"net/url"
)

type HeaderValue []string
type Headers map[string][]string

type RequestLine struct {
	URI     url.URL
	Version string
	Method  common.HttpMethod
}

func (h Headers) Set(key, value string) {
	canonicalKey := common.GetCanonicalName(key)

	h[canonicalKey] = []string{value}
}

// It performs the upsert operation where if the key does not exist it will create a new entry else it will append to existing values.
func (h Headers) Upsert(key, value string) {
	canonicalKey := common.GetCanonicalName(key)

	existingValues, exists := h[canonicalKey]

	// If the key does not exist set it.
	if !exists {
		h.Set(key, value)
	}

	h[canonicalKey] = append(existingValues, value)
}

// Returns the value of the header. It fetches the key by canonical name which it automatically converts to when getting it.
func (h Headers) Get(key string) (value HeaderValue, exist bool) {
	canonicalKey := common.GetCanonicalName(key)
	value, exist = h[canonicalKey]

	return
}

func (h Headers) Remove(key string) {
	canonicalKey := common.GetCanonicalName(key)
	delete(h, canonicalKey)
}
