package threatcrowd

import (
	"encoding/json"
	"fmt"

	"github.com/drsigned/sigsubs/pkg/sources"
)

type response struct {
	Subdomains []string `json:"subdomains"`
}

// Source is the passive scraping agent
type Source struct{}

// Run function returns all subdomains found with the service
func (source *Source) Run(domain string, session *sources.Session) chan sources.Subdomain {
	subdomains := make(chan sources.Subdomain)

	go func() {
		defer close(subdomains)

		res, _ := session.SimpleGet(fmt.Sprintf("https://www.threatcrowd.org/searchApi/v2/domain/report/?domain=%s", domain))

		var results response

		if err := json.Unmarshal(res.Body(), &results); err != nil {
			return
		}

		for _, i := range results.Subdomains {
			subdomains <- sources.Subdomain{Source: source.Name(), Value: i}
		}
	}()

	return subdomains
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "threatcrowd"
}
