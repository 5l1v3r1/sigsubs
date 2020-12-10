package sublist3r

import (
	"encoding/json"
	"fmt"

	"github.com/drsigned/sigsubs/pkg/sources"
)

// Source is the passive sources agent
type Source struct{}

// Run function returns all subdomains found with the service
func (source *Source) Run(domain string, session *sources.Session) chan sources.Subdomain {
	subdomains := make(chan sources.Subdomain)

	go func() {
		defer close(subdomains)

		res, _ := session.SimpleGet(fmt.Sprintf("https://api.sublist3r.com/search.php?domain=%s", domain))

		body := res.Body()

		var results []string

		if err := json.Unmarshal(body, &results); err != nil {
			return
		}

		for _, i := range results {
			subdomains <- sources.Subdomain{Source: source.Name(), Value: i}
		}
	}()

	return subdomains
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "sublist3r"
}
