package sources

import (
	"regexp"
	"sync"
)

var extractorMutex = &sync.Mutex{}

// NewExtractor creates a new regular expression to extract
// subdomains from text based on the given domain.
func NewExtractor(domain string) (*regexp.Regexp, error) {
	extractorMutex.Lock()
	defer extractorMutex.Unlock()

	extractor, err := regexp.Compile(`[a-zA-Z0-9\*_.-]+\.` + domain)
	if err != nil {
		return nil, err
	}

	return extractor, nil
}
