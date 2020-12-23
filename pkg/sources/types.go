package sources

import (
	"regexp"

	"github.com/valyala/fasthttp"
)

// Subdomain is a result structure returned by a source
type Subdomain struct {
	Source string
	Value  string
}

// Source is an interface inherited by each passive source
type Source interface {
	// Run takes a domain as argument and a session object
	// which contains the extractor for subdomains, http client
	// and other stuff.
	Run(string, *Session) chan Subdomain
	// Name returns the name of the source
	Name() string
}

// Keys contains the current API Keys we have in store
type Keys struct {
	Chaos  string   `json:"chaos"`
	GitHub []string `json:"github"`
}

// Session is the option passed to the source, an option is created
// uniquely for eac source.
type Session struct {
	// Extractor is the regex for subdomains created for each domain
	Extractor *regexp.Regexp
	// Keys is the API keys for the application
	Keys *Keys
	// Client is the current http client
	Client *fasthttp.Client
}
