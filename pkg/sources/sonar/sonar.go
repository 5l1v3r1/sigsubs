package sonar

import (
	"encoding/json"
	"fmt"

	"github.com/drsigned/sigsubs/pkg/sources"
)

type response struct {
	Subdomains []string
}

// Source is the passive sources agent
type Source struct{}

// Run function returns all subdomains found with the service
func (source *Source) Run(domain string, session *sources.Session) chan sources.Subdomain {
	subdomains := make(chan sources.Subdomain)

	go func() {
		defer close(subdomains)

		res, _ := session.SimpleGet(fmt.Sprintf("https://sonar.omnisint.io/subdomains/%s", domain))

		var results []interface{}

		if err := json.Unmarshal(res.Body(), &results); err != nil {
			return
		}

		for _, i := range results {
			subdomains <- sources.Subdomain{Source: source.Name(), Value: i.(string)}
		}
	}()

	return subdomains
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "sonar"
}
